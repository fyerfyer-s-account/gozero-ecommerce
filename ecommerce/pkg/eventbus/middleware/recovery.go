package middleware

import (
    "context"
    "fmt"
    "runtime/debug"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

func Recovery(next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, msg amqp.Delivery) (err error) {
        defer func() {
            if r := recover(); r != nil {
                stack := debug.Stack()
                
                fields := map[string]interface{}{
                    "panic": r,
                    "stack": string(stack),
                }

                zerolog.GetLogger().Error(ctx, "Recovered from panic in message handler", fmt.Errorf("%v", r), fields)

                err = fmt.Errorf("panic recovered: %v", r)
                msg.Nack(false, true) // Requeue message
            }
        }()

        return next(ctx, msg)
    }
}