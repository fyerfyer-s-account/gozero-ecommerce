package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/consumer/retry"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"github.com/streadway/amqp"
)

type EventHandler interface {
	Handle(event *types.CartEvent) error
}

type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	config     *config.RabbitMQConfig
	handlers   map[types.EventType][]EventHandler
	retrier    *retry.Retrier
	middleware []middleware.Middleware
	logger     middleware.Logger
}

func NewConsumer(config *config.RabbitMQConfig, logger middleware.Logger) (*Consumer, error) {
	conn, err := amqp.Dial(config.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	backoff := retry.NewExponentialBackoff(
		time.Duration(config.Retry.InitialInterval)*time.Millisecond,
		time.Duration(config.Retry.MaxInterval)*time.Millisecond,
		config.Retry.BackoffFactor,
		config.Retry.Jitter,
	)
	retrier := retry.NewRetrier(config.Retry.MaxAttempts, backoff)

	c := &Consumer{
		conn:     conn,
		channel:  channel,
		config:   config,
		handlers: make(map[types.EventType][]EventHandler),
		retrier:  retrier,
		logger:   logger,
	}

	if config.Middleware.EnableRecovery {
		c.Use(middleware.NewRecoveryMiddleware(logger))
	}
	if config.Middleware.EnableLogging {
		c.Use(middleware.NewLoggerMiddleware(logger))
	}

	return c, nil
}

func (c *Consumer) Use(m middleware.Middleware) {
	c.middleware = append(c.middleware, m)
}

func (c *Consumer) Subscribe(eventType types.EventType, handler EventHandler) {
	c.handlers[eventType] = append(c.handlers[eventType], handler)
}

func (c *Consumer) Start() error {
	err := c.setupExchangesAndQueues()
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		c.config.Queues.CartStatus.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go c.handleMessages(msgs)
	return nil
}

func (c *Consumer) setupExchangesAndQueues() error {
	err := c.channel.ExchangeDeclare(
		c.config.Exchanges.CartEvent.Name,
		c.config.Exchanges.CartEvent.Type,
		c.config.Exchanges.CartEvent.Durable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = c.channel.QueueDeclare(
		c.config.Queues.CartStatus.Name,
		c.config.Queues.CartStatus.Durable,
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = c.channel.QueueBind(
		c.config.Queues.CartStatus.Name,
		c.config.Queues.CartStatus.RoutingKey,
		c.config.Exchanges.CartEvent.Name,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	return nil
}

func (c *Consumer) handleMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		handler := c.createHandlerChain(func(msg amqp.Delivery) error {
			var event types.CartEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				return &types.RetryableError{Err: err}
			}

			handlers, exists := c.handlers[event.Type]
			if !exists {
				return nil
			}

			for _, h := range handlers {
				if err := c.retrier.DoWithContext(context.Background(), func() error {
					return h.Handle(&event)
				}); err != nil {
					return err
				}
			}
			return nil
		})

		if err := handler(msg); err != nil {
			c.handleFailedMessage(msg, err)
		}
	}
}

func (c *Consumer) createHandlerChain(handler middleware.HandlerFunc) middleware.HandlerFunc {
	chain := handler
	for i := len(c.middleware) - 1; i >= 0; i-- {
		chain = c.middleware[i](chain)
	}
	return chain
}

func (c *Consumer) handleFailedMessage(msg amqp.Delivery, err error) {
	if err := c.channel.Publish(
		c.config.DeadLetter.Exchange,
		c.config.DeadLetter.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: msg.ContentType,
			Body:        msg.Body,
			MessageId:   msg.MessageId,
			Timestamp:   time.Now(),
			Headers:     map[string]interface{}{"error": err.Error()},
		},
	); err != nil {
		c.logger.Error("failed to publish to dead letter exchange",
			"error", err,
			"message_id", msg.MessageId,
		)
	}
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
