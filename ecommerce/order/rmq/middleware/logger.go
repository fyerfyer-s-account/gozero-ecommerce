package middleware

import (
	"time"

	"github.com/streadway/amqp"
)

// HandlerFunc defines the message handler function type
type HandlerFunc func(msg amqp.Delivery) error

// Middleware defines the middleware function type
type Middleware func(HandlerFunc) HandlerFunc

// Logger is an interface for logging
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// LoggerMiddleware is a middleware for logging
type LoggerMiddleware struct {
	logger Logger
}

// NewLoggerMiddleware creates a new LoggerMiddleware
func NewLoggerMiddleware(logger Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(msg amqp.Delivery) error {
			start := time.Now()

			logger.Info("processing message",
				"message_id", msg.MessageId,
				"routing_key", msg.RoutingKey,
			)

			err := next(msg)

			duration := time.Since(start)

			if err != nil {
				logger.Error("failed to process message",
					"message_id", msg.MessageId,
					"routing_key", msg.RoutingKey,
					"error", err,
					"duration", duration,
				)
				return err
			}

			logger.Info("message processed successfully",
				"message_id", msg.MessageId,
				"routing_key", msg.RoutingKey,
				"duration", duration,
			)

			return nil
		}
	}
}
