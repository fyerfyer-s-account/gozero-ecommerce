package middleware

import (
    "time"
    "github.com/streadway/amqp"
)

type HandlerFunc func(msg amqp.Delivery) error

type Middleware func(HandlerFunc) HandlerFunc

type Logger interface {
    Info(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
}

func NewLoggerMiddleware(logger Logger) Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(msg amqp.Delivery) error {
            start := time.Now()

            logger.Info("processing cart message",
                "message_id", msg.MessageId,
                "routing_key", msg.RoutingKey,
            )

            err := next(msg)

            duration := time.Since(start)

            if err != nil {
                logger.Error("failed to process cart message",
                    "message_id", msg.MessageId,
                    "routing_key", msg.RoutingKey,
                    "error", err,
                    "duration_ms", duration.Milliseconds(),
                )
                return err
            }

            logger.Info("cart message processed successfully",
                "message_id", msg.MessageId,
                "routing_key", msg.RoutingKey,
                "duration_ms", duration.Milliseconds(),
            )

            return nil
        }
    }
}