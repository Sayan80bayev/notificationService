package messaging

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	queue      string
	routingKey string
	handler    func(eventType string, data []byte)
	logger     *logrus.Logger
}

// NewRabbitConsumer creates a new consumer bound to an exchange and routing key.
func NewRabbitConsumer(amqpURL, exchange, queue, routingKey string, handler func(eventType string, data []byte), logger *logrus.Logger) (*RabbitConsumer, error) {
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

	return &RabbitConsumer{
		conn:       conn,
		channel:    ch,
		exchange:   exchange,
		queue:      q.Name,
		routingKey: routingKey,
		handler:    handler,
		logger:     logger,
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
			c.handler(msg.RoutingKey, msg.Body)
		}
	}
}

func (c *RabbitConsumer) Close() {
	_ = c.channel.Close()
	_ = c.conn.Close()
	c.logger.Info("RabbitMQ consumer closed")
}
