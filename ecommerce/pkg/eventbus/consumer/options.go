package consumer

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/streadway/amqp"
)

type Options struct {
    Exchange     string
    Queue        string
    RoutingKey   string
    AutoAck      bool
    Middlewares  []middleware.MiddlewareFunc
    QueueDeclare bool
    QueueOptions amqp.Table
}

type Option func(*Options)

func WithExchange(exchange string) Option {
    return func(o *Options) {
        o.Exchange = exchange
    }
}

func WithQueue(queue string) Option {
    return func(o *Options) {
        o.Queue = queue
    }
}

func WithRoutingKey(key string) Option {
    return func(o *Options) {
        o.RoutingKey = key
    }
}

func WithAutoAck(autoAck bool) Option {
    return func(o *Options) {
        o.AutoAck = autoAck
    }
}

func WithMiddlewares(mws ...middleware.MiddlewareFunc) Option {
    return func(o *Options) {
        o.Middlewares = append(o.Middlewares, mws...)
    }
}

func WithQueueDeclare(declare bool) Option {
    return func(o *Options) {
        o.QueueDeclare = declare
    }
}

func WithQueueOptions(options amqp.Table) Option {
    return func(o *Options) {
        o.QueueOptions = options
    }
}