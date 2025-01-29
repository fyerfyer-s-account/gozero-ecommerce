package producer

import (
    "context"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type PaymentProducer struct {
    logger    *zerolog.Logger
    producer  *producer.Producer
    exchange  string
}

func NewPaymentProducer(ch *amqp.Channel, exchange string) *PaymentProducer {
    p := producer.NewProducer(ch,
        producer.WithExchange(exchange),
        producer.WithPersistentDelivery(),
        producer.WithAsync(),
        producer.WithBatch(100, time.Second),
    )

    return &PaymentProducer{
        logger:   zerolog.GetLogger(),
        producer: p,
        exchange: exchange,
    }
}

func (p *PaymentProducer) PublishPaymentCreated(ctx context.Context, event *types.PaymentCreatedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *PaymentProducer) PublishPaymentSuccess(ctx context.Context, event *types.PaymentSuccessEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *PaymentProducer) PublishPaymentFailed(ctx context.Context, event *types.PaymentFailedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *PaymentProducer) PublishPaymentRefund(ctx context.Context, event *types.PaymentRefundEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *PaymentProducer) PublishPaymentVerification(ctx context.Context, event *types.PaymentVerificationEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *PaymentProducer) Close() error {
    return p.producer.Close()
}