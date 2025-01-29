package svc

import (
    "context"
    "fmt"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config         config.Config
    CartItemsModel model.CartItemsModel
    CartStatsModel model.CartStatisticsModel
    ProductRpc     productservice.ProductService
    Producer       *producer.CartProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
    // Initialize MySQL connection
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

    // Establish RabbitMQ connection
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := rmqBroker.Connect(ctx); err != nil {
        panic(fmt.Sprintf("Failed to connect to RabbitMQ broker: %v", err))
    }

    // Get RabbitMQ channel
    ch, err := rmqBroker.Channel()
    if err != nil {
        panic(fmt.Sprintf("Failed to create RabbitMQ channel: %v", err))
    }

    return &ServiceContext{
        Config:         c,
        CartItemsModel: model.NewCartItemsModel(conn, c.CacheRedis),
        CartStatsModel: model.NewCartStatisticsModel(conn, c.CacheRedis),
        ProductRpc:     productservice.NewProductService(zrpc.MustNewClient(c.ProductRpc)),
        Producer:       producer.NewCartProducer(ch, c.RabbitMQ.Exchange),
    }
}