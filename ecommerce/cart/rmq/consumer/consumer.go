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
    logger               *zerolog.Logger
    channel             *amqp.Channel
    statusHandler       *handlers.StatusHandler
    orderHandler        *handlers.OrderHandler
    inventoryHandler    *handlers.InventoryHandler
    clearHandler        *handlers.ClearHandler
    selectionHandler    *handlers.SelectionHandler
    paymentSuccessHandler *handlers.PaymentSuccessHandler
    consumers           []*consumer.Consumer
}

func NewCartConsumer(
    ch *amqp.Channel,
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *CartConsumer {
    return &CartConsumer{
        logger:  zerolog.GetLogger(),
        channel: ch,
        statusHandler: handlers.NewStatusHandler(
            cartItemsModel,
            cartStatsModel,
        ),
        orderHandler: handlers.NewOrderHandler(
            cartItemsModel,
            cartStatsModel,
        ),
        inventoryHandler: handlers.NewInventoryHandler(
            cartItemsModel,
            cartStatsModel,
        ),
        clearHandler: handlers.NewClearHandler(
            cartItemsModel,
            cartStatsModel,
        ),
        selectionHandler: handlers.NewSelectionHandler(
            cartItemsModel,
            cartStatsModel,
        ),
        paymentSuccessHandler: handlers.NewPaymentSuccessHandler(
            cartItemsModel,
            cartStatsModel,
        ),
    }
}

func (c *CartConsumer) Start(ctx context.Context) error {
    // Create handlers with their respective configurations
    consumers := []struct {
        name        string
        handler     middleware.HandlerFunc
        queue       string
        exchange    string
        routingKey  string
    }{
        {
            name:       "status",
            handler:    c.statusHandler.Handle,
            queue:      "cart.status",
            exchange:   "cart.events",
            routingKey: "cart.status.*",
        },
        {
            name:       "order",
            handler:    c.orderHandler.Handle,
            queue:      "cart.order",
            exchange:   "order.events",
            routingKey: "order.*",
        },
        {
            name:       "inventory",
            handler:    c.inventoryHandler.Handle,
            queue:      "cart.inventory",
            exchange:   "inventory.events",
            routingKey: "inventory.stock.*",
        },
        {
            name:       "clear",
            handler:    c.clearHandler.Handle,
            queue:      "cart.clear",
            exchange:   "cart.events",
            routingKey: "cart.clear",
        },
        {
            name:       "selection",
            handler:    c.selectionHandler.Handle,
            queue:      "cart.selection",
            exchange:   "cart.events",
            routingKey: "cart.select.*",
        },
        {
            name:       "payment.success",
            handler:    c.paymentSuccessHandler.Handle,
            queue:      "cart.payment.success",
            exchange:   "payment.events",
            routingKey: "payment.success",
        },
    }

    // Initialize and start all consumers
    c.consumers = make([]*consumer.Consumer, 0, len(consumers))
    for _, cfg := range consumers {
        cons, err := consumer.NewConsumer(
            c.channel,
            cfg.handler,
            consumer.WithQueue(cfg.queue),
            consumer.WithExchange(cfg.exchange),
            consumer.WithRoutingKey(cfg.routingKey),
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

        if err := cons.Start(ctx); err != nil {
            return err
        }

        c.consumers = append(c.consumers, cons)
    }

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