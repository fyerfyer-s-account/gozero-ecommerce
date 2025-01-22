package middleware

import (
    "fmt"
    "github.com/streadway/amqp"
    "runtime/debug"
)

type RecoveryMiddleware struct {
    logger Logger
}

func NewRecoveryMiddleware(logger Logger) Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(msg amqp.Delivery) (err error) {
            defer func() {
                if r := recover(); r != nil {
                    stack := debug.Stack()
                    err = fmt.Errorf("panic recovered in inventory processing: %v\n%s", r, stack)
                    logger.Error("panic recovered in inventory message processing",
                        "error", err,
                        "message_id", msg.MessageId,
                        "stack", string(stack),
                    )
                }
            }()
            
            return next(msg)
        }
    }
}