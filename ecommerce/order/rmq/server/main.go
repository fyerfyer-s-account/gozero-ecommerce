package main

import (
	"context"
	"flag"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "../rpc/etc/order.yaml", "the config file")

func main() {
    flag.Parse()

    var c config.Config
    conf.MustLoad(*configFile, &c)

    logger := &zerolog.LogWrapper{Logger: logx.WithContext(context.TODO())}
    
    cons, err := consumer.NewConsumer(&c.RabbitMQ, logger)
    if err != nil {
        panic(err)
    }
    defer cons.Close()

    if err := cons.Start(); err != nil {
        panic(err)
    }

    logx.Info("Order RMQ server started")
    select {}
}