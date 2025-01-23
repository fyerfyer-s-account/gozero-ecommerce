package main

import (
	"context"
	"flag"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

var configFile = flag.String("f", "../rpc/etc/inventory.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logger := &zerolog.LogWrapper{Logger: logx.WithContext(context.TODO())}

	cons, err := consumer.NewConsumer(&c.RabbitMQ, logger, 
		inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)), 
		messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc)))
	if err != nil {
		panic(err)
	}
	defer cons.Close()

	if err := cons.Start(); err != nil {
		panic(err)
	}

	logx.Info("Inventory RMQ server started")
	select {}
}
