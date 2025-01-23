package svc

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config            config.Config
    StocksModel       model.StocksModel
    StockLocksModel   model.StockLocksModel
    StockRecordsModel model.StockRecordsModel
    WarehousesModel   model.WarehousesModel
    Producer          *producer.Producer
    MessageRpc        messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)

    // Initialize RPC clients
    messageRpc := messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc))

    // Convert config to RabbitMQ config
    rmqConfig := &rmqconfig.RabbitMQConfig{
        Host:     c.RabbitMQ.Host,
        Port:     c.RabbitMQ.Port,
        Username: c.RabbitMQ.Username,
        Password: c.RabbitMQ.Password,
        VHost:    c.RabbitMQ.VHost,
        Exchanges: rmqconfig.ExchangeConfigs{
            InventoryEvent: rmqconfig.ExchangeConfig{
                Name:    c.RabbitMQ.Exchanges.InventoryEvent.Name,
                Type:    c.RabbitMQ.Exchanges.InventoryEvent.Type,
                Durable: c.RabbitMQ.Exchanges.InventoryEvent.Durable,
            },
        },
        Queues: rmqconfig.QueueConfigs{
            StockUpdate: rmqconfig.QueueConfig{
                Name:       c.RabbitMQ.Queues.StockUpdate.Name,
                RoutingKey: c.RabbitMQ.Queues.StockUpdate.RoutingKey,
                Durable:    c.RabbitMQ.Queues.StockUpdate.Durable,
            },
            StockAlert: rmqconfig.QueueConfig{
                Name:       c.RabbitMQ.Queues.StockAlert.Name,
                RoutingKey: c.RabbitMQ.Queues.StockAlert.RoutingKey,
                Durable:    c.RabbitMQ.Queues.StockAlert.Durable,
            },
            StockLock: rmqconfig.QueueConfig{
                Name:       c.RabbitMQ.Queues.StockLock.Name,
                RoutingKey: c.RabbitMQ.Queues.StockLock.RoutingKey,
                Durable:    c.RabbitMQ.Queues.StockLock.Durable,
            },
            OrderEvents: rmqconfig.QueueConfig{
                Name:       c.RabbitMQ.Queues.OrderEvents.Name,
                RoutingKey: c.RabbitMQ.Queues.OrderEvents.RoutingKey,
                Durable:    c.RabbitMQ.Queues.OrderEvents.Durable,
            },
        },
    }

    // Initialize producer
    prod, err := producer.NewProducer(rmqConfig)
    if err != nil {
        panic(err)
    }

    return &ServiceContext{
        Config:            c,
        StocksModel:       model.NewStocksModel(conn, c.CacheRedis),
        StockLocksModel:   model.NewStockLocksModel(conn, c.CacheRedis),
        StockRecordsModel: model.NewStockRecordsModel(conn, c.CacheRedis),
        WarehousesModel:   model.NewWarehousesModel(conn, c.CacheRedis),
        Producer:          prod,
        MessageRpc:        messageRpc,
    }
}