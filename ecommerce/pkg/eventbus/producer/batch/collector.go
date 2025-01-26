package batch

import (
    "context"
    "github.com/streadway/amqp"
    "time"
)

type Collector struct {
    batchSize    int
    batchTimeout time.Duration
    messages     []*message
    msgChan      chan *message
    done         chan struct{}
    sender       *Sender
}

type message struct {
    exchange    string
    routingKey  string
    publishing  amqp.Publishing
}

func NewCollector(size int, timeout time.Duration, sender *Sender) *Collector {
    c := &Collector{
        batchSize:    size,
        batchTimeout: timeout,
        messages:     make([]*message, 0, size),
        msgChan:      make(chan *message, size),
        done:         make(chan struct{}),
        sender:       sender,
    }
    go c.start()
    return c
}

func (c *Collector) Collect(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case c.msgChan <- &message{exchange, routingKey, msg}:
        return nil
    }
}

func (c *Collector) start() {
    timer := time.NewTimer(c.batchTimeout)
    defer timer.Stop()

    for {
        select {
        case <-c.done:
            return
        case msg := <-c.msgChan:
            c.messages = append(c.messages, msg)
            if len(c.messages) >= c.batchSize {
                c.flush()
                timer.Reset(c.batchTimeout)
            }
        case <-timer.C:
            if len(c.messages) > 0 {
                c.flush()
            }
            timer.Reset(c.batchTimeout)
        }
    }
}

func (c *Collector) flush() {
    c.sender.Send(c.messages)
    c.messages = c.messages[:0]
}

func (c *Collector) Close() {
    close(c.done)
}