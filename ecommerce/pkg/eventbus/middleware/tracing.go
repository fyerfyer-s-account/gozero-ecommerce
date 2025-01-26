package middleware

import (
    "context"

    "github.com/streadway/amqp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func Tracing(serviceName string) MiddlewareFunc {
    return func(next HandlerFunc) HandlerFunc {
        return func(ctx context.Context, msg amqp.Delivery) error {
            tracer := otel.Tracer(serviceName)

            ctx, span := tracer.Start(ctx, msg.RoutingKey,
                trace.WithAttributes(
                    attribute.String("messaging.system", "rabbitmq"),
                    attribute.String("messaging.destination", msg.Exchange),
                    attribute.String("messaging.routing_key", msg.RoutingKey),
                    attribute.String("messaging.message_id", msg.MessageId),
                ),
            )
            defer span.End()

            if err := next(ctx, msg); err != nil {
                span.RecordError(err)
                return err
            }

            return nil
        }
    }
}