package main

import (
    "context"
    "flag"
    "fmt"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/server"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/svc"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
    flag.Parse()

    var c config.Config
    conf.MustLoad(*configFile, &c)

    ctx := svc.NewServiceContext(c)
    srv := server.NewRmqServer(context.Background(), ctx)

    serviceGroup := service.NewServiceGroup()
    serviceGroup.Add(srv)

    fmt.Printf("Starting order rmq server at %s...\n", c.ListenOn)
    serviceGroup.Start()
}