package consumer

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type CartConsumer struct {
    logger           *zerolog.Logger
    channel          *amqp.Channel
    statusHandler    *handlers.StatusHandler
    clearHandler     *handlers.ClearHandler
    selectionHandler *handlers.SelectionHandler
    inventoryHandler *handlers.InventoryHandler
    orderHandler     *handlers.OrderHandler
    consumers        []*consumer.Consumer
}

func NewCartConsumer(
    ch *amqp.Channel,
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *CartConsumer {
    return &CartConsumer{
        logger:           zerolog.GetLogger(),
        channel:          ch,
        statusHandler:    handlers.NewStatusHandler(cartItemsModel, cartStatsModel),
        clearHandler:     handlers.NewClearHandler(cartItemsModel, cartStatsModel),
        selectionHandler: handlers.NewSelectionHandler(cartItemsModel, cartStatsModel),
        inventoryHandler: handlers.NewInventoryHandler(cartItemsModel, cartStatsModel),
        orderHandler:     handlers.NewOrderHandler(cartItemsModel, cartStatsModel),
    }
}

func (c *CartConsumer) Start(ctx context.Context) error {
    // Create status consumer
    statusConsumer, err := consumer.NewConsumer(
        c.channel,
        c.statusHandler.Handle,
        consumer.WithQueue("cart.status"),
        consumer.WithExchange("cart.events"),
        consumer.WithRoutingKey("cart.status.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("cart-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create clear consumer
    clearConsumer, err := consumer.NewConsumer(
        c.channel,
        c.clearHandler.Handle,
        consumer.WithQueue("cart.clear"),
        consumer.WithExchange("cart.events"),
        consumer.WithRoutingKey("cart.clear.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("cart-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create selection consumer
    selectionConsumer, err := consumer.NewConsumer(
        c.channel,
        c.selectionHandler.Handle,
        consumer.WithQueue("cart.selection"),
        consumer.WithExchange("cart.events"),
        consumer.WithRoutingKey("cart.selection.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("cart-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create inventory consumer
    inventoryConsumer, err := consumer.NewConsumer(
        c.channel,
        c.inventoryHandler.Handle,
        consumer.WithQueue("cart.inventory"),
        consumer.WithExchange("inventory.events"),
        consumer.WithRoutingKey("inventory.stock.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("cart-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create order consumer
    orderConsumer, err := consumer.NewConsumer(
        c.channel,
        c.orderHandler.Handle,
        consumer.WithQueue("cart.order"),
        consumer.WithExchange("order.events"),
        consumer.WithRoutingKey("order.status.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("cart-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Start all consumers
    if err := statusConsumer.Start(ctx); err != nil {
        return err
    }
    if err := clearConsumer.Start(ctx); err != nil {
        return err
    }
    if err := selectionConsumer.Start(ctx); err != nil {
        return err
    }
    if err := inventoryConsumer.Start(ctx); err != nil {
        return err
    }
    if err := orderConsumer.Start(ctx); err != nil {
        return err
    }

    c.consumers = append(c.consumers, 
        statusConsumer, 
        clearConsumer, 
        selectionConsumer,
        inventoryConsumer,
        orderConsumer,
    )
    return nil
}

func (c *CartConsumer) Stop() error {
    for _, consumer := range c.consumers {
        if err := consumer.Stop(); err != nil {
            c.logger.Error(context.Background(), "Failed to stop consumer", err, nil)
        }
    }
    return nil
}