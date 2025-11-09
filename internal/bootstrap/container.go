package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/caching"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/Sayan80bayev/go-project/pkg/messaging"
	_ "github.com/lib/pq"
	"notificationService/internal/config"
	ms "notificationService/internal/messaging"
	"notificationService/internal/repository"
	"notificationService/internal/service"
	"time"
)

// Container holds all dependencies
type Container struct {
	DB                     *sql.DB
	Redis                  caching.CacheService
	Producer               messaging.Producer
	Consumer               messaging.Consumer
	Config                 *config.Config
	NotificationService    service.NotificationService
	NotificationRepository repository.NotificationRepository
	JWKSUrl                string
}

// Init initializes all dependencies and returns a container
func Init() (*Container, error) {
	logger := logging.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := initPostgresDatabase(cfg)
	if err != nil {
		return nil, err
	}

	cacheService, err := initRedis(cfg)
	if err != nil {
		return nil, err
	}

	consumer, err := initRabbitMQConsumer(cfg)
	if err != nil {
		return nil, err
	}

	jwksURL := buildJWKSURL(cfg)

	nr := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(nr)

	logger.Info("✅ Dependencies initialized successfully")

	return &Container{
		DB:                     db,
		Redis:                  cacheService,
		Consumer:               consumer,
		NotificationService:    svc,
		NotificationRepository: nr,
		Config:                 cfg,
		JWKSUrl:                jwksURL,
	}, nil
}

// initPostgresDatabase initializes a PostgreSQL database connection
func initPostgresDatabase(cfg *config.Config) (*sql.DB, error) {
	logger := logging.GetLogger()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Ping the database to verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	logger.Info("PostgreSQL connected")
	return db, nil
}

func initRedis(cfg *config.Config) (*caching.RedisService, error) {
	logger := logging.GetLogger()
	redisCache, err := caching.NewRedisService(caching.RedisConfig{
		DB:       0,
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
	})

	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Info("Redis connected")
	return redisCache, nil
}

func initRabbitMQConsumer(cfg *config.Config) (messaging.Consumer, error) {
	amqpUrl := buildAmqpURL(cfg)
	consumer, err := ms.NewRabbitConsumer(amqpUrl, cfg.RabbitMQExchange, cfg.RabbitMQQueue, cfg.RabbitMQRoutingKey, nil, logging.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("kafka consumer init failed: %w", err)
	}

	logging.GetLogger().Infof("Kafka consumer initialized")
	return consumer, nil
}

func buildJWKSURL(cfg *config.Config) string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", cfg.KeycloakURL, cfg.KeycloakRealm)
}

func buildAmqpURL(cfg *config.Config) string {
	return fmt.Sprintf("ampq://%s:%s/%s:%s/", cfg.RabbitMQUser, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort)
}
