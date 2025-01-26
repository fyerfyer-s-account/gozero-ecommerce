package batch

import (
    "github.com/streadway/amqp"
)

type Sender struct {
    ch *amqp.Channel
}

func NewSender(ch *amqp.Channel) *Sender {
    return &Sender{ch: ch}
}

func (s *Sender) Send(messages []*message) {
    for _, msg := range messages {
        s.ch.Publish(
            msg.exchange,
            msg.routingKey,
            false,
            false,
            msg.publishing,
        )
    }
}