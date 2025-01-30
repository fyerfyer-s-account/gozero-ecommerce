package consumer

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type MarketingConsumer struct {
    logger                *zerolog.Logger
    channel              *amqp.Channel
    couponHandler        *handlers.CouponEventHandler
    promotionHandler     *handlers.PromotionEventHandler
    pointsHandler        *handlers.PointsEventHandler
    paymentSuccessHandler *handlers.PaymentSuccessHandler
    paymentFailedHandler  *handlers.PaymentFailedHandler
    consumers            []*consumer.Consumer
}

func NewMarketingConsumer(
    ch *amqp.Channel,
    couponsModel model.CouponsModel,
    userCouponsModel model.UserCouponsModel,
    promotionsModel model.PromotionsModel,
    userPointsModel model.UserPointsModel,
    pointsRecordModel model.PointsRecordsModel,
) *MarketingConsumer {
    return &MarketingConsumer{
        logger:  zerolog.GetLogger(),
        channel: ch,
        couponHandler: handlers.NewCouponEventHandler(
            couponsModel,
            userCouponsModel,
        ),
        promotionHandler: handlers.NewPromotionEventHandler(
            promotionsModel,
        ),
        pointsHandler: handlers.NewPointsEventHandler(
            userPointsModel,
            pointsRecordModel,
        ),
        paymentSuccessHandler: handlers.NewPaymentSuccessHandler(
            userPointsModel,
            userCouponsModel,
            pointsRecordModel,
        ),
        paymentFailedHandler: handlers.NewPaymentFailedHandler(
            userCouponsModel,
            userPointsModel,
            pointsRecordModel,
        ),
    }
}

func (c *MarketingConsumer) Start(ctx context.Context) error {
    // Create coupon event consumer
    couponConsumer, err := consumer.NewConsumer(
        c.channel,
        c.couponHandler.Handle,
        consumer.WithQueue("marketing.coupon"),
        consumer.WithExchange("marketing.events"),
        consumer.WithRoutingKey("marketing.coupon.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("marketing-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create promotion event consumer
    promotionConsumer, err := consumer.NewConsumer(
        c.channel,
        c.promotionHandler.Handle,
        consumer.WithQueue("marketing.promotion"),
        consumer.WithExchange("marketing.events"),
        consumer.WithRoutingKey("marketing.promotion.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("marketing-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create points event consumer
    pointsConsumer, err := consumer.NewConsumer(
        c.channel,
        c.pointsHandler.Handle,
        consumer.WithQueue("marketing.points"),
        consumer.WithExchange("marketing.events"),
        consumer.WithRoutingKey("marketing.points.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("marketing-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create payment success consumer
    paymentSuccessConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentSuccessHandler.Handle,
        consumer.WithQueue("marketing.payment.success"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.success"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("marketing-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create payment failed consumer
    paymentFailedConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentFailedHandler.Handle,
        consumer.WithQueue("marketing.payment.failed"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.failed"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("marketing-consumer"),
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
        {"coupon", couponConsumer},
        {"promotion", promotionConsumer},
        {"points", pointsConsumer},
        {"payment.success", paymentSuccessConsumer},
        {"payment.failed", paymentFailedConsumer},
    } {
        if err := c.consumer.Start(ctx); err != nil {
            return err
        }
    }

    c.consumers = []*consumer.Consumer{
        couponConsumer,
        promotionConsumer,
        pointsConsumer,
        paymentSuccessConsumer,
        paymentFailedConsumer,
    }

    return nil
}

func (c *MarketingConsumer) Stop() error {
    for _, consumer := range c.consumers {
        if err := consumer.Stop(); err != nil {
            c.logger.Error(context.Background(), "Failed to stop consumer", err, nil)
        }
    }
    return nil
}