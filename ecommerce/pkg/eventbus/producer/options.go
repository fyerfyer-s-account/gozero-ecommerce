package producer

import (
    "github.com/streadway/amqp"
    "time"
)

type Options struct {
    Exchange      string
    RoutingKey    string
    Mandatory     bool
    Immediate     bool
    Headers       amqp.Table
    DeliveryMode  uint8
    Async         bool
    BatchSize     int
    BatchTimeout  time.Duration
}

type Option func(*Options)

func WithExchange(exchange string) Option {
    return func(o *Options) {
        o.Exchange = exchange
    }
}

func WithRoutingKey(key string) Option {
    return func(o *Options) {
        o.RoutingKey = key
    }
}

func WithHeaders(headers amqp.Table) Option {
    return func(o *Options) {
        o.Headers = headers
    }
}

func WithPersistentDelivery() Option {
    return func(o *Options) {
        o.DeliveryMode = amqp.Persistent
    }
}

func WithAsync() Option {
    return func(o *Options) {
        o.Async = true
    }
}

func WithBatch(size int, timeout time.Duration) Option {
    return func(o *Options) {
        o.BatchSize = size
        o.BatchTimeout = timeout
    }
}