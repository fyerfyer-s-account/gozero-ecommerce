package async

import (
    "context"
    "github.com/streadway/amqp"
)

type Dispatcher struct {
    channel chan *message
    done    chan struct{}
    ch      *amqp.Channel
}

type message struct {
    exchange    string
    routingKey  string
    publishing  amqp.Publishing
}

func NewDispatcher(ch *amqp.Channel) *Dispatcher {
    d := &Dispatcher{
        channel: make(chan *message, 1000),
        done:    make(chan struct{}),
        ch:      ch,
    }
    go d.start()
    return d
}

func (d *Dispatcher) Dispatch(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case d.channel <- &message{exchange, routingKey, msg}:
        return nil
    }
}

func (d *Dispatcher) start() {
    for {
        select {
        case <-d.done:
            return
        case msg := <-d.channel:
            d.ch.Publish(
                msg.exchange,
                msg.routingKey,
                false,
                false,
                msg.publishing,
            )
        }
    }
}

func (d *Dispatcher) Close() {
    close(d.done)
}