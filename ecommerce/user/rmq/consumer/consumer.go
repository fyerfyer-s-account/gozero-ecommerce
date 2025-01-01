package consumer

import (
	"context"
	"encoding/json"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rmq/types"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	config   *config.RabbitMQConfig
	handlers map[types.EventType]EventHandler
}

type EventHandler func(event *types.UserEvent) error

func NewConsumer(config *config.RabbitMQConfig) (*Consumer, error) {
	conn, err := amqp.Dial(config.GetDSN())
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:     conn,
		channel:  ch,
		config:   config,
		handlers: make(map[types.EventType]EventHandler),
	}, nil
}

func (c *Consumer) RegisterHandler(eventType types.EventType, handler EventHandler) {
	c.handlers[eventType] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.config.Queues.UserNotification.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event types.UserEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				msg.Nack(false, true)
				continue
			}

			if handler, ok := c.handlers[event.Type]; ok {
				if err := handler(&event); err != nil {
					msg.Nack(false, true)
					continue
				}
			}

			msg.Ack(false)
		}
	}()

	<-ctx.Done()
	return nil
}
