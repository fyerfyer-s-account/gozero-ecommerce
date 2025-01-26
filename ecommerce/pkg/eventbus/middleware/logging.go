package middleware

import (
    "context"
    "time"

    "github.com/streadway/amqp"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
)

func Logging(next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, msg amqp.Delivery) error {
        start := time.Now()
        logger := zerolog.GetLogger()

        fields := map[string]interface{}{
            "exchange": msg.Exchange,
            "routing_key": msg.RoutingKey,
            "message_id": msg.MessageId,
        }

        logger.Info(ctx, "Processing message started", fields)

        err := next(ctx, msg)

        fields["duration"] = time.Since(start)
        fields["success"] = err == nil
        
        if err != nil {
            logger.WithError(ctx, err, "Processing message completed", fields)
        } else {
            logger.Info(ctx, "Processing message completed", fields)
        }

        return err
    }
}