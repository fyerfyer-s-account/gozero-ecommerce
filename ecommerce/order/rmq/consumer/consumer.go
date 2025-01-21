package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/consumer/retry"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
	"github.com/streadway/amqp"
	"time"
)

type EventHandler interface {
	Handle(event *types.OrderEvent) error
}

type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	config     *config.RabbitMQConfig
	handlers   map[types.EventType][]EventHandler
	retrier    *retry.Retrier
	middleware []middleware.Middleware  // Changed from HandlerFunc to Middleware
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

	// Initialize retrier
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

	// Setup middleware
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
		c.config.Queues.OrderStatus.Name,
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
	// Declare exchange
	err := c.channel.ExchangeDeclare(
		c.config.Exchanges.OrderEvent.Name,
		c.config.Exchanges.OrderEvent.Type,
		c.config.Exchanges.OrderEvent.Durable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Declare queues and bind them
	for _, queueConfig := range []config.QueueConfig{
		c.config.Queues.OrderStatus,
		c.config.Queues.OrderAlert,
	} {
		_, err = c.channel.QueueDeclare(
			queueConfig.Name,
			queueConfig.Durable,
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %v", queueConfig.Name, err)
		}

		err = c.channel.QueueBind(
			queueConfig.Name,
			queueConfig.RoutingKey,
			c.config.Exchanges.OrderEvent.Name,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s: %v", queueConfig.Name, err)
		}
	}

	return nil
}

func (c *Consumer) handleMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		// Create handler chain
		handler := c.createHandlerChain(func(msg amqp.Delivery) error {
			var event types.OrderEvent
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
	// Publish to dead letter exchange
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
