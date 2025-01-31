package producer

import (
    "context"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type MessageProducer struct {
    logger    *zerolog.Logger
    producer  *producer.Producer
    exchange  string
}

func NewMessageProducer(ch *amqp.Channel, exchange string) *MessageProducer {
    p := producer.NewProducer(ch,
        producer.WithExchange(exchange),
        producer.WithPersistentDelivery(),
        producer.WithAsync(),
        producer.WithBatch(100, time.Second),
    )

    return &MessageProducer{
        logger:   zerolog.GetLogger(),
        producer: p,
        exchange: exchange,
    }
}

func (p *MessageProducer) PublishMessageEvent(ctx context.Context, event *types.MessageEventSentEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MessageProducer) PublishTemplateEvent(ctx context.Context, event *types.MessageTemplateEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MessageProducer) PublishBatchEvent(ctx context.Context, event *types.MessageBatchEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MessageProducer) PublishPaymentSuccessEvent(ctx context.Context, event *types.MessagePaymentSuccessEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MessageProducer) PublishPaymentFailedEvent(ctx context.Context, event *types.MessagePaymentFailedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MessageProducer) Close() error {
    return p.producer.Close()
}