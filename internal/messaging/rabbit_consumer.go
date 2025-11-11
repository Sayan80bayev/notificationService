package messaging

import (
	"context"
	"fmt"
	"github.com/Sayan80bayev/go-project/pkg/logging"
	"github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}
type RabbitConsumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	queue      string
	routingKey string
	handlers   map[string]EventHandler
	logger     *logrus.Logger
}

type EventHandler func(data []byte) error

func NewRabbitConsumer(amqpURL, exchange, queue, routingKey string, logger *logrus.Logger) (*RabbitConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	logger.Infof("RabbitMQ consumer connected (queue=%s, routing=%s)", q.Name, routingKey)
	handlers := make(map[string]EventHandler)
	return &RabbitConsumer{
		conn:       conn,
		channel:    ch,
		exchange:   exchange,
		queue:      q.Name,
		routingKey: routingKey,
		logger:     logger,
		handlers:   handlers,
	}, nil
}

func (c *RabbitConsumer) Start(ctx context.Context) {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Errorf("failed to register consumer: %v", err)
		return
	}

	c.logger.Infof("[Consumer] Listening on queue=%s routing=%s", c.queue, c.routingKey)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("[Consumer] Context canceled, shutting down...")
			return
		case msg := <-msgs:
			if msg.Body == nil {
				continue
			}
			c.logger.Infof("[Consumer] Received event=%s body=%s", msg.RoutingKey, string(msg.Body))
			c.handleRabbitEvent(msg.RoutingKey, msg.Body)
		}
	}
}

func (c *RabbitConsumer) handleRabbitEvent(eventType string, data []byte) {
	logger := logging.GetLogger()

	if handler, ok := c.handlers[eventType]; ok {
		if err := handler(data); err != nil {
			logger.Warnf("Failed to handle event (%s): %v", eventType, err)
			return
		}
	} else {
		logger.Warnf("Unrecognized event type: %s (payload: %s)", eventType, string(data))
	}
}

func (c *RabbitConsumer) RegisterHandler(eventType string, handler EventHandler) {
	c.handlers[eventType] = handler
}

func (c *RabbitConsumer) Close() {
	_ = c.channel.Close()
	_ = c.conn.Close()
	c.logger.Info("RabbitMQ consumer closed")
}
