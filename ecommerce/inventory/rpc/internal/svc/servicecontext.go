package svc

import (
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config            config.Config
    StocksModel       model.StocksModel
    StockLocksModel   model.StockLocksModel
    StockRecordsModel model.StockRecordsModel
    WarehousesModel   model.WarehousesModel
    Producer          *producer.InventoryProducer
    MessageRpc        messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)

    // Initialize RabbitMQ broker
    rmqBroker := broker.NewAMQPBroker(&rmqconfig.RabbitMQConfig{
        Host:              c.RabbitMQ.Host,
        Port:              c.RabbitMQ.Port,
        Username:          c.RabbitMQ.Username,
        Password:          c.RabbitMQ.Password,
        VHost:             c.RabbitMQ.VHost,
        ConnectionTimeout: time.Duration(c.RabbitMQ.ConnectionTimeout) * time.Second,
        HeartbeatInterval: time.Duration(c.RabbitMQ.HeartbeatInterval) * time.Second,
    })

    ch, err := rmqBroker.Channel()
    if err != nil {
        panic(err)
    }

    return &ServiceContext{
        Config:            c,
        StocksModel:       model.NewStocksModel(conn, c.CacheRedis),
        StockLocksModel:   model.NewStockLocksModel(conn, c.CacheRedis),
        StockRecordsModel: model.NewStockRecordsModel(conn, c.CacheRedis),
        WarehousesModel:   model.NewWarehousesModel(conn, c.CacheRedis),
        Producer:          producer.NewInventoryProducer(ch, c.RabbitMQ.Exchange),
        MessageRpc:        messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc)),
    }
}