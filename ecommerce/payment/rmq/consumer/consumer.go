package consumer

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type PaymentConsumer struct {
    logger               *zerolog.Logger
    channel              *amqp.Channel
    orderPaymentHandler  *handlers.OrderPaymentHandler
    orderRefundHandler   *handlers.OrderRefundHandler
    verificationHandler  *handlers.PaymentVerificationHandler
    consumers            []*consumer.Consumer
}

func NewPaymentConsumer(
    ch *amqp.Channel,
    paymentOrders model.PaymentOrdersModel,
    paymentChannels model.PaymentChannelsModel,
    refundOrders model.RefundOrdersModel,
    paymentLogs model.PaymentLogsModel,
) *PaymentConsumer {
    return &PaymentConsumer{
        logger:  zerolog.GetLogger(),
        channel: ch,
        orderPaymentHandler: handlers.NewOrderPaymentHandler(
            paymentOrders,
            paymentChannels,
            paymentLogs,
        ),
        orderRefundHandler: handlers.NewOrderRefundHandler(
            paymentOrders,
            refundOrders,
            paymentLogs,
        ),
        verificationHandler: handlers.NewPaymentVerificationHandler(
            paymentOrders,
            paymentChannels,
            paymentLogs,
        ),
    }
}

func (c *PaymentConsumer) Start(ctx context.Context) error {
    // Create order payment consumer
    paymentConsumer, err := consumer.NewConsumer(
        c.channel,
        c.orderPaymentHandler.Handle,
        consumer.WithQueue("payment.order"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("order.payment.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("payment-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create refund consumer
    refundConsumer, err := consumer.NewConsumer(
        c.channel,
        c.orderRefundHandler.Handle,
        consumer.WithQueue("payment.refund"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("order.refund.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("payment-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create verification consumer
    verificationConsumer, err := consumer.NewConsumer(
        c.channel,
        c.verificationHandler.Handle,
        consumer.WithQueue("payment.verification"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.verification.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("payment-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Start consumers
    if err := paymentConsumer.Start(ctx); err != nil {
        return err
    }
    if err := refundConsumer.Start(ctx); err != nil {
        return err
    }
    if err := verificationConsumer.Start(ctx); err != nil {
        return err
    }

    c.consumers = append(c.consumers, paymentConsumer, refundConsumer, verificationConsumer)
    return nil
}

func (c *PaymentConsumer) Stop() error {
    for _, consumer := range c.consumers {
        if err := consumer.Stop(); err != nil {
            c.logger.Error(context.Background(), "Failed to stop consumer", err, nil)
        }
    }
    return nil
}