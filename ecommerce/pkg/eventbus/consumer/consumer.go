package consumer

import (
    "context"
    "fmt"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type Consumer struct {
    channel *amqp.Channel
    options Options
    logger  *zerolog.Logger
    handler middleware.HandlerFunc
}

func NewConsumer(ch *amqp.Channel, handler middleware.HandlerFunc, opts ...Option) (*Consumer, error) {
    options := Options{
        AutoAck: false,
        QueueDeclare: true,
        QueueOptions: amqp.Table{},
    }

    for _, opt := range opts {
        opt(&options)
    }

    // Apply default middlewares
    defaultMiddlewares := []middleware.MiddlewareFunc{
        middleware.Recovery,
        middleware.Logging,
        middleware.Retry(nil),
    }
    options.Middlewares = append(defaultMiddlewares, options.Middlewares...)

    // Chain middlewares
    finalHandler := handler
    for i := len(options.Middlewares) - 1; i >= 0; i-- {
        finalHandler = options.Middlewares[i](finalHandler)
    }

    return &Consumer{
        channel: ch,
        options: options,
        logger:  zerolog.GetLogger(),
        handler: finalHandler,
    }, nil
}

func (c *Consumer) Start(ctx context.Context) error {
    if c.options.QueueDeclare {
        if _, err := c.channel.QueueDeclare(
            c.options.Queue,
            true,  // durable
            false, // auto-delete
            false, // exclusive
            false, // no-wait
            c.options.QueueOptions,
        ); err != nil {
            return fmt.Errorf("failed to declare queue: %w", err)
        }

        if err := c.channel.QueueBind(
            c.options.Queue,
            c.options.RoutingKey,
            c.options.Exchange,
            false,
            nil,
        ); err != nil {
            return fmt.Errorf("failed to bind queue: %w", err)
        }
    }

    msgs, err := c.channel.Consume(
        c.options.Queue,
        "",    // consumer
        c.options.AutoAck,
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return fmt.Errorf("failed to start consuming: %w", err)
    }

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-msgs:
                if !ok {
                    return
                }
                if err := c.handler(ctx, msg); err != nil {
                    c.logger.Error(ctx, "Failed to process message", err, map[string]interface{}{
                        "exchange": msg.Exchange,
                        "routing_key": msg.RoutingKey,
                        "message_id": msg.MessageId,
                    })
                    if !c.options.AutoAck {
                        msg.Nack(false, true)
                    }
                } else if !c.options.AutoAck {
                    msg.Ack(false)
                }
            }
        }
    }()

    return nil
}

func (c *Consumer) Stop() error {
    if err := c.channel.Close(); err != nil {
        return fmt.Errorf("failed to close channel: %w", err)
    }
    return nil
}