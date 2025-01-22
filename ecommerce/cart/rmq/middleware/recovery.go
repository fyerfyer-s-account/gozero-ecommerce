package middleware

import (
    "fmt"
    "github.com/streadway/amqp"
    "runtime/debug"
)

func NewRecoveryMiddleware(logger Logger) Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(msg amqp.Delivery) (err error) {
            defer func() {
                if r := recover(); r != nil {
                    stack := debug.Stack()
                    err = fmt.Errorf("panic recovered in cart message processing: %v\n%s", r, stack)

                    logger.Error("panic recovered in cart message processing",
                        "error", r,
                        "message_id", msg.MessageId,
                        "stack", string(stack),
                    )

                    // Ensure message is not lost
                    msg.Nack(false, true)
                }
            }()

            return next(msg)
        }
    }
}