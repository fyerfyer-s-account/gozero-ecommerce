package producer

import (
	"context"
	"encoding/json"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rmq/types"
	"github.com/streadway/amqp"
	"sync"
)

type Producer struct {
	config   *config.RabbitMQConfig
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	mu       sync.RWMutex
}

func NewProducer(cfg *config.RabbitMQConfig) (*Producer, error) {
	p := &Producer{
		config:   cfg,
		exchange: cfg.Exchanges.ProductEvent.Name,
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Producer) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := amqp.Dial(p.config.GetDSN())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		p.exchange,
		p.config.Exchanges.ProductEvent.Type,
		p.config.Exchanges.ProductEvent.Durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	p.conn = conn
	p.channel = ch
	return nil
}

func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

func (p *Producer) PublishEvent(ctx context.Context, event *types.ProductEvent) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		p.exchange,
		string(event.Type),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}
