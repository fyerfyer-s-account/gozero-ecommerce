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

type LoggerMiddleware struct {
    logger Logger
}

func NewLoggerMiddleware(logger Logger) Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(msg amqp.Delivery) error {
            start := time.Now()

            logger.Info("processing inventory message",
                "message_id", msg.MessageId,
                "routing_key", msg.RoutingKey,
            )

            err := next(msg)

            duration := time.Since(start)

            if err != nil {
                logger.Error("failed to process inventory message",
                    "message_id", msg.MessageId,
                    "routing_key", msg.RoutingKey,
                    "error", err,
                    "duration", duration,
                )
                return err
            }

            logger.Info("inventory message processed successfully",
                "message_id", msg.MessageId,
                "routing_key", msg.RoutingKey,
                "duration", duration,
            )

            return nil
        }
    }
}