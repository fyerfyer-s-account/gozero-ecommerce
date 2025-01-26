package producer

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer/async"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer/batch"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/serializer/json"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type Producer struct {
    channel     *amqp.Channel
    options     Options
    logger      *zerolog.Logger
    serializer  *json.JsonSerializer
    dispatcher  *async.Dispatcher
    collector   *batch.Collector
}

func NewProducer(ch *amqp.Channel, opts ...Option) *Producer {
    options := Options{
        DeliveryMode: amqp.Persistent,
        Headers:      make(amqp.Table),
    }

    for _, opt := range opts {
        opt(&options)
    }

    p := &Producer{
        channel:    ch,
        options:    options,
        logger:     zerolog.GetLogger(),
        serializer: json.New(),
    }

    if options.Async {
        p.dispatcher = async.NewDispatcher(ch)
    }

    if options.BatchSize > 0 {
        p.collector = batch.NewCollector(options.BatchSize, options.BatchTimeout, 
            batch.NewSender(ch))
    }

    return p
}

func (p *Producer) Publish(ctx context.Context, msg interface{}) error {
    data, err := p.serializer.Marshal(msg)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }

    publishing := amqp.Publishing{
        ContentType:  p.serializer.ContentType(),
        Body:        data,
        DeliveryMode: p.options.DeliveryMode,
        Headers:     p.options.Headers,
    }

    if p.options.Async {
        return p.dispatcher.Dispatch(ctx, p.options.Exchange, p.options.RoutingKey, publishing)
    }

    if p.collector != nil {
        return p.collector.Collect(ctx, p.options.Exchange, p.options.RoutingKey, publishing)
    }

    return p.channel.Publish(
        p.options.Exchange,
        p.options.RoutingKey,
        p.options.Mandatory,
        p.options.Immediate,
        publishing,
    )
}

func (p *Producer) Close() error {
    if p.dispatcher != nil {
        p.dispatcher.Close()
    }
    if p.collector != nil {
        p.collector.Close()
    }
    return p.channel.Close()
}