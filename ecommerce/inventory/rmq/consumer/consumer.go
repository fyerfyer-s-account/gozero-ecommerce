package consumer

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type InventoryConsumer struct {
    logger        *zerolog.Logger
    channel       *amqp.Channel
    updateHandler *handlers.UpdateHandler
    alertHandler  *handlers.AlertHandler
    lockHandler   *handlers.LockHandler
    orderHandler  *handlers.OrderHandler
    consumers     []*consumer.Consumer
}

func NewInventoryConsumer(
    ch *amqp.Channel,
    stocksModel model.StocksModel,
    stockLocksModel model.StockLocksModel,
    stockRecordsModel model.StockRecordsModel,
) *InventoryConsumer {
    return &InventoryConsumer{
        logger:  zerolog.GetLogger(),
        channel: ch,
        updateHandler: handlers.NewUpdateHandler(
            stocksModel,
            stockRecordsModel,
        ),
        alertHandler: handlers.NewAlertHandler(
            stocksModel,
            stockRecordsModel,
        ),
        lockHandler: handlers.NewLockHandler(
            stocksModel,
            stockLocksModel,
            stockRecordsModel,
        ),
        orderHandler: handlers.NewOrderHandler(
            stocksModel,
            stockLocksModel,
            stockRecordsModel,
        ),
    }
}

func (c *InventoryConsumer) Start(ctx context.Context) error {
    // Create stock update consumer
    updateConsumer, err := consumer.NewConsumer(
        c.channel,
        c.updateHandler.Handle,
        consumer.WithQueue("stock.update"),
        consumer.WithExchange("inventory.events"),
        consumer.WithRoutingKey("inventory.stock.update.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("inventory-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create stock alert consumer
    alertConsumer, err := consumer.NewConsumer(
        c.channel,
        c.alertHandler.Handle,
        consumer.WithQueue("stock.alert"),
        consumer.WithExchange("inventory.events"),
        consumer.WithRoutingKey("inventory.stock.alert.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("inventory-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create stock lock consumer
    lockConsumer, err := consumer.NewConsumer(
        c.channel,
        c.lockHandler.Handle,
        consumer.WithQueue("stock.lock"),
        consumer.WithExchange("inventory.events"),
        consumer.WithRoutingKey("inventory.stock.lock.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("inventory-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create order events consumer
    orderConsumer, err := consumer.NewConsumer(
        c.channel,
        c.orderHandler.Handle,
        consumer.WithQueue("order.events"),
        consumer.WithExchange("order.events"),
        consumer.WithRoutingKey("order.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("inventory-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Start all consumers
    if err := updateConsumer.Start(ctx); err != nil {
        return err
    }
    if err := alertConsumer.Start(ctx); err != nil {
        return err
    }
    if err := lockConsumer.Start(ctx); err != nil {
        return err
    }
    if err := orderConsumer.Start(ctx); err != nil {
        return err
    }

    c.consumers = append(c.consumers, updateConsumer, alertConsumer, lockConsumer, orderConsumer)
    return nil
}

func (c *InventoryConsumer) Stop() error {
    for _, consumer := range c.consumers {
        if err := consumer.Stop(); err != nil {
            c.logger.Error(context.Background(), "Failed to stop consumer", err, nil)
        }
    }
    return nil
}