package consumer

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer/retry"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
    "github.com/streadway/amqp"
    "time"
)

type EventHandler interface {
    Handle(event *types.InventoryEvent) error
}

type Consumer struct {
    conn          *amqp.Connection
    channel       *amqp.Channel
    config        *config.RabbitMQConfig
    handlers      map[types.EventType][]EventHandler
    retrier       *retry.Retrier
    middleware    []middleware.Middleware
    logger        middleware.Logger
    inventoryRpc  inventoryclient.Inventory
    messageRpc    messageservice.MessageService
}

func NewConsumer(
    config *config.RabbitMQConfig,
    logger middleware.Logger,
    inventoryRpc inventoryclient.Inventory,
    messageRpc messageservice.MessageService,
) (*Consumer, error) {
    conn, err := amqp.Dial(config.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open channel: %v", err)
    }

    backoff := retry.NewExponentialBackoff(
        time.Duration(config.Retry.InitialInterval)*time.Millisecond,
        time.Duration(config.Retry.MaxInterval)*time.Millisecond,
        config.Retry.BackoffFactor,
        config.Retry.Jitter,
    )
    retrier := retry.NewRetrier(config.Retry.MaxAttempts, backoff)

    c := &Consumer{
        conn:          conn,
        channel:       channel,
        config:        config,
        handlers:      make(map[types.EventType][]EventHandler),
        retrier:       retrier,
        logger:        logger,
        inventoryRpc:  inventoryRpc,
        messageRpc:    messageRpc,
    }

    if config.Middleware.EnableRecovery {
        c.Use(middleware.NewRecoveryMiddleware(logger))
    }
    if config.Middleware.EnableLogging {
        c.Use(middleware.NewLoggerMiddleware(logger))
    }

    return c, nil
}

func (c *Consumer) Use(m middleware.Middleware) {
    c.middleware = append(c.middleware, m)
}

func (c *Consumer) Subscribe(eventType types.EventType, handler EventHandler) {
    c.handlers[eventType] = append(c.handlers[eventType], handler)
}

func (c *Consumer) Start() error {
    err := c.setupExchangesAndQueues()
    if err != nil {
        return err
    }

    for _, queueConfig := range []config.QueueConfig{
        c.config.Queues.StockUpdate,
        c.config.Queues.StockAlert,
        c.config.Queues.StockLock,
    } {
        msgs, err := c.channel.Consume(
            queueConfig.Name,
            "",    // consumer
            false, // auto-ack
            false, // exclusive
            false, // no-local
            false, // no-wait
            nil,   // args
        )
        if err != nil {
            return fmt.Errorf("failed to register a consumer for queue %s: %v", queueConfig.Name, err)
        }

        go c.handleMessages(msgs)
    }

    return nil
}

func (c *Consumer) setupExchangesAndQueues() error {
    err := c.channel.ExchangeDeclare(
        c.config.Exchanges.InventoryEvent.Name,
        c.config.Exchanges.InventoryEvent.Type,
        c.config.Exchanges.InventoryEvent.Durable,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    for _, queueConfig := range []config.QueueConfig{
        c.config.Queues.StockUpdate,
        c.config.Queues.StockAlert,
        c.config.Queues.StockLock,
    } {
        _, err = c.channel.QueueDeclare(
            queueConfig.Name,
            queueConfig.Durable,
            false,
            false,
            false,
            nil,
        )
        if err != nil {
            return fmt.Errorf("failed to declare queue %s: %v", queueConfig.Name, err)
        }

        err = c.channel.QueueBind(
            queueConfig.Name,
            queueConfig.RoutingKey,
            c.config.Exchanges.InventoryEvent.Name,
            false,
            nil,
        )
        if err != nil {
            return fmt.Errorf("failed to bind queue %s: %v", queueConfig.Name, err)
        }
    }

    return nil
}

func (c *Consumer) handleMessages(msgs <-chan amqp.Delivery) {
    for msg := range msgs {
        handler := c.createHandlerChain(func(msg amqp.Delivery) error {
            var event types.InventoryEvent
            if err := json.Unmarshal(msg.Body, &event); err != nil {
                return &types.RetryableError{Err: err}
            }

            handlers, exists := c.handlers[event.Type]
            if !exists {
                return nil
            }

            for _, h := range handlers {
                if err := c.retrier.DoWithContext(context.Background(), func() error {
                    return h.Handle(&event)
                }); err != nil {
                    return err
                }
            }
            return nil
        })

        if err := handler(msg); err != nil {
            c.handleFailedMessage(msg, err)
            msg.Nack(false, false)
        } else {
            msg.Ack(false)
        }
    }
}

func (c *Consumer) createHandlerChain(handler middleware.HandlerFunc) middleware.HandlerFunc {
    chain := handler
    for i := len(c.middleware) - 1; i >= 0; i-- {
        chain = c.middleware[i](chain)
    }
    return chain
}

func (c *Consumer) handleFailedMessage(msg amqp.Delivery, err error) {
    if err := c.channel.Publish(
        c.config.DeadLetter.Exchange,
        c.config.DeadLetter.RoutingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: msg.ContentType,
            Body:       msg.Body,
            MessageId:  msg.MessageId,
            Timestamp:  time.Now(),
            Headers:    map[string]interface{}{"error": err.Error()},
        },
    ); err != nil {
        c.logger.Error("failed to publish to dead letter exchange",
            "error", err,
            "message_id", msg.MessageId,
        )
    }
}

func (c *Consumer) Close() error {
    if err := c.channel.Close(); err != nil {
        return err
    }
    return c.conn.Close()
}