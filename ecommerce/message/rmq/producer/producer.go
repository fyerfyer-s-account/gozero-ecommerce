package producer

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
    "github.com/google/uuid"
    "github.com/streadway/amqp"
)

type Producer struct {
    conn         *amqp.Connection
    channel      *amqp.Channel
    config       config.RabbitMQConfig
    mutex        sync.RWMutex
    notifyClose  chan *amqp.Error
}

func NewProducer(cfg config.RabbitMQConfig) *Producer {
    return &Producer{
        config: cfg,
    }
}

func (p *Producer) Start() error {
    return p.connect()
}

func (p *Producer) Stop() {
    if p.channel != nil {
        p.channel.Close()
    }
    if p.conn != nil {
        p.conn.Close()
    }
}

func (p *Producer) connect() error {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    var err error
    p.conn, err = amqp.Dial(p.config.GetDSN())
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    p.channel, err = p.conn.Channel()
    if err != nil {
        return fmt.Errorf("failed to open channel: %v", err)
    }

    err = p.channel.ExchangeDeclare(
        p.config.Exchanges.MessageEvent.Name,
        p.config.Exchanges.MessageEvent.Type,
        p.config.Exchanges.MessageEvent.Durable,
        false, // auto-delete
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    p.notifyClose = make(chan *amqp.Error)
    p.conn.NotifyClose(p.notifyClose)

    go p.handleReconnect()

    return nil
}

func (p *Producer) handleReconnect() {
    for err := range p.notifyClose {
        log.Printf("Connection closed, reconnecting: %v", err)
        for {
            if err := p.connect(); err == nil {
                log.Println("Reconnected to RabbitMQ")
                break
            }
            time.Sleep(5 * time.Second)
        }
    }
}

func (p *Producer) PublishEvent(ctx context.Context, eventType types.EventType, routingKey string, data interface{}, metadata types.Metadata) error {
    event := &types.MessageEvent{
        ID:        uuid.New().String(),
        Type:      eventType,
        Timestamp: time.Now(),
        Data:      data,
        Metadata:  metadata,
    }

    body, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %v", err)
    }

    p.mutex.RLock()
    defer p.mutex.RUnlock()

    err = p.channel.Publish(
        p.config.Exchanges.MessageEvent.Name,
        routingKey,
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:        body,
            MessageId:   event.ID,
            Timestamp:   event.Timestamp,
            DeliveryMode: 2, // persistent
        },
    )

    if err != nil {
        return fmt.Errorf("failed to publish message: %v", err)
    }

    return nil
}

// Helper methods for specific event types

func (p *Producer) PublishMessageCreated(ctx context.Context, data *types.MessageCreatedData, metadata types.Metadata) error {
    return p.PublishEvent(ctx, types.EventTypeMessageCreated, "message.created", data, metadata)
}

func (p *Producer) PublishMessageSent(ctx context.Context, data *types.MessageSentData, metadata types.Metadata) error {
    return p.PublishEvent(ctx, types.EventTypeMessageSent, "message.sent", data, metadata)
}

func (p *Producer) PublishMessageRead(ctx context.Context, data *types.MessageReadData, metadata types.Metadata) error {
    return p.PublishEvent(ctx, types.EventTypeMessageRead, "message.read", data, metadata)
}

func (p *Producer) PublishTemplateCreated(ctx context.Context, data *types.TemplateData, metadata types.Metadata) error {
    return p.PublishEvent(ctx, types.EventTypeTemplateCreated, "template.created", data, metadata)
}

func (p *Producer) PublishTemplateUpdated(ctx context.Context, data *types.TemplateData, metadata types.Metadata) error {
    return p.PublishEvent(ctx, types.EventTypeTemplateUpdated, "template.updated", data, metadata)
}