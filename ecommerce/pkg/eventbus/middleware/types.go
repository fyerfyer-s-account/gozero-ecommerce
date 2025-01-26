package middleware

import (
    "context"
    
    "github.com/streadway/amqp"
)

type HandlerFunc func(context.Context, amqp.Delivery) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc