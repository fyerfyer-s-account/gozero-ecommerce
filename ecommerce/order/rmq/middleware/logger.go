package middleware

import (
    "time"
    "github.com/streadway/amqp"
)

type Logger interface {
    Info(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
}

type LoggerMiddleware struct {
    logger Logger
}

func NewLoggerMiddleware(logger Logger) *LoggerMiddleware {
    return &LoggerMiddleware{
        logger: logger,
    }
}

func (m *LoggerMiddleware) Handle(next HandlerFunc) HandlerFunc {
    return func(msg amqp.Delivery) error {
        start := time.Now()
        
        m.logger.Info("processing message",
            "message_id", msg.MessageId,
            "routing_key", msg.RoutingKey,
        )
        
        err := next(msg)
        
        duration := time.Since(start)
        
        if err != nil {
            m.logger.Error("failed to process message",
                "message_id", msg.MessageId,
                "routing_key", msg.RoutingKey,
                "error", err,
                "duration", duration,
            )
            return err
        }
        
        m.logger.Info("message processed successfully",
            "message_id", msg.MessageId,
            "routing_key", msg.RoutingKey,
            "duration", duration,
        )
        
        return nil
    }
}

type HandlerFunc func(amqp.Delivery) error