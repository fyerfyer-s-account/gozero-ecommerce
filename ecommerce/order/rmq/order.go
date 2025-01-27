package main

import (
	"context"
	"flag"
	"log"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/server"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
    log.Println("start order rmq service!!!")
    flag.Parse()

    log.Println("parse yaml file!!!")
    var c config.Config
    conf.MustLoad(*configFile, &c)

    log.Println("setup servicecontext!!!")
    ctx := svc.NewServiceContext(c)

    log.Println("setup rmq server!!!")
    srv := server.NewRmqServer(context.Background(), ctx)

    log.Println("set up servicegroup!!!")
    serviceGroup := service.NewServiceGroup()
    serviceGroup.Add(srv)

    log.Printf("Starting order rmq server at %s...\n", c.ListenOn)
    serviceGroup.Start()
}