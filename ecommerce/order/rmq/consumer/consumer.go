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
    logger               *zerolog.Logger
    channel             *amqp.Channel
    statusHandler       *handlers.StatusHandler
    alertHandler        *handlers.AlertHandler
    paymentSuccessHandler *handlers.PaymentSuccessHandler
    paymentFailedHandler  *handlers.PaymentFailedHandler
    paymentRefundHandler  *handlers.PaymentRefundHandler
    consumers           []*consumer.Consumer
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
        paymentSuccessHandler: handlers.NewPaymentSuccessHandler(
            ordersModel,
            paymentsModel,
        ),
        paymentFailedHandler: handlers.NewPaymentFailedHandler(
            ordersModel,
            paymentsModel,
        ),
        paymentRefundHandler: handlers.NewPaymentRefundHandler(
            ordersModel,
            paymentsModel,
            refundsModel,
        ),
    }
}

func (c *OrderConsumer) Start(ctx context.Context) error {
    // Create existing consumers
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

    // Create payment success consumer
    paymentSuccessConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentSuccessHandler.Handle,
        consumer.WithQueue("order.payment.success"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.success"),
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

    // Create payment failed consumer
    paymentFailedConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentFailedHandler.Handle,
        consumer.WithQueue("order.payment.failed"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.failed"),
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

    // Create payment refund consumer
    paymentRefundConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentRefundHandler.Handle,
        consumer.WithQueue("order.payment.refund"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.refund"),
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

    // Start all consumers
    for _, c := range []struct {
        name     string
        consumer *consumer.Consumer
    }{
        {"status", statusConsumer},
        {"alert", alertConsumer},
        {"payment.success", paymentSuccessConsumer},
        {"payment.failed", paymentFailedConsumer},
        {"payment.refund", paymentRefundConsumer},
    } {
        if err := c.consumer.Start(ctx); err != nil {
            return err
        }
    }

    c.consumers = []*consumer.Consumer{
        statusConsumer,
        alertConsumer,
        paymentSuccessConsumer,
        paymentFailedConsumer,
        paymentRefundConsumer,
    }

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