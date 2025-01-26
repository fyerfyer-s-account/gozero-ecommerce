package consumer

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type OrderConsumer struct {
    logger        *zerolog.Logger
    channel       *amqp.Channel
    statusHandler *handlers.StatusHandler
    alertHandler  *handlers.AlertHandler
    consumers     []*consumer.Consumer
}

func NewOrderConsumer(
    ch *amqp.Channel,
    ordersModel model.OrdersModel,
    paymentsModel model.OrderPaymentsModel,
    shippingModel model.OrderShippingModel,
    refundsModel model.OrderRefundsModel,
) *OrderConsumer {
    return &OrderConsumer{
        logger:  zerolog.GetLogger(),
        channel: ch,
        statusHandler: handlers.NewStatusHandler(
            ordersModel,
            paymentsModel,
            shippingModel,
            refundsModel,
        ),
        alertHandler: handlers.NewAlertHandler(
            ordersModel,
            paymentsModel,
            refundsModel,
        ),
    }
}

func (c *OrderConsumer) Start(ctx context.Context) error {
    // Create status consumer
    statusConsumer, err := consumer.NewConsumer(
        c.channel,
        c.statusHandler.Handle,
        consumer.WithQueue("order.status"),
        consumer.WithExchange("order.events"),
        consumer.WithRoutingKey("order.status.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("order-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create alert consumer
    alertConsumer, err := consumer.NewConsumer(
        c.channel,
        c.alertHandler.Handle,
        consumer.WithQueue("order.alert"),
        consumer.WithExchange("order.events"),
        consumer.WithRoutingKey("order.alert.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("order-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Start consumers
    if err := statusConsumer.Start(ctx); err != nil {
        return err
    }
    if err := alertConsumer.Start(ctx); err != nil {
        return err
    }

    c.consumers = append(c.consumers, statusConsumer, alertConsumer)
    return nil
}

func (c *OrderConsumer) Stop() error {
    for _, consumer := range c.consumers {
        if err := consumer.Stop(); err != nil {
            c.logger.Error(context.Background(), "Failed to stop consumer", err, nil)
        }
    }
    return nil
}