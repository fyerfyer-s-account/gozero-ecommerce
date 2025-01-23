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

var configFile = flag.String("f", "../etc/inventory.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logger := &zerolog.LogWrapper{Logger: logx.WithContext(context.TODO())}

	// Convert config to RabbitMQ config
    rmqConfig := &config.RabbitMQConfig{
        Host:     c.RabbitMQ.Host,
        Port:     c.RabbitMQ.Port,
        Username: c.RabbitMQ.Username,
        Password: c.RabbitMQ.Password,
        VHost:    c.RabbitMQ.VHost,
        Exchanges: config.ExchangeConfigs{
            InventoryEvent: config.ExchangeConfig{
                Name:    c.RabbitMQ.Exchanges.InventoryEvent.Name,
                Type:    c.RabbitMQ.Exchanges.InventoryEvent.Type,
                Durable: c.RabbitMQ.Exchanges.InventoryEvent.Durable,
            },
        },
        Queues: config.QueueConfigs{
            StockUpdate: config.QueueConfig{
                Name:       c.RabbitMQ.Queues.StockUpdate.Name,
                RoutingKey: c.RabbitMQ.Queues.StockUpdate.RoutingKey,
                Durable:    c.RabbitMQ.Queues.StockUpdate.Durable,
            },
            StockAlert: config.QueueConfig{
                Name:       c.RabbitMQ.Queues.StockAlert.Name,
                RoutingKey: c.RabbitMQ.Queues.StockAlert.RoutingKey,
                Durable:    c.RabbitMQ.Queues.StockAlert.Durable,
            },
            StockLock: config.QueueConfig{
                Name:       c.RabbitMQ.Queues.StockLock.Name,
                RoutingKey: c.RabbitMQ.Queues.StockLock.RoutingKey,
                Durable:    c.RabbitMQ.Queues.StockLock.Durable,
            },
            OrderEvents: config.QueueConfig{
                Name:       c.RabbitMQ.Queues.OrderEvents.Name,
                RoutingKey: c.RabbitMQ.Queues.OrderEvents.RoutingKey,
                Durable:    c.RabbitMQ.Queues.OrderEvents.Durable,
            },
        },
    }

	cons, err := consumer.NewConsumer(rmqConfig, logger, 
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
