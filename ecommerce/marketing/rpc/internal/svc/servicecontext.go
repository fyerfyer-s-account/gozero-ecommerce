package svc

import (
    "context"
    "fmt"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config             config.Config
    CouponsModel       model.CouponsModel
    UserCouponsModel   model.UserCouponsModel
    PromotionsModel    model.PromotionsModel
    UserPointsModel    model.UserPointsModel
    PointsRecordsModel model.PointsRecordsModel
    Producer           *producer.MarketingProducer
    MessageRpc         messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)

    // Initialize RabbitMQ broker
    rmqConfig := &rmqconfig.RabbitMQConfig{
        Host:              c.RabbitMQ.Host,
        Port:              c.RabbitMQ.Port,
        Username:          c.RabbitMQ.Username,
        Password:          c.RabbitMQ.Password,
        VHost:             c.RabbitMQ.VHost,
        ConnectionTimeout: time.Duration(c.RabbitMQ.ConnectionTimeout) * time.Second,
        HeartbeatInterval: time.Duration(c.RabbitMQ.HeartbeatInterval) * time.Second,
    }

    rmqBroker := broker.NewAMQPBroker(rmqConfig)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := rmqBroker.Connect(ctx); err != nil {
        panic(fmt.Sprintf("Failed to connect to RabbitMQ broker: %v", err))
    }

    ch, err := rmqBroker.Channel()
    if err != nil {
        panic(fmt.Sprintf("Failed to create RabbitMQ channel: %v", err))
    }

    return &ServiceContext{
        Config:             c,
        CouponsModel:       model.NewCouponsModel(conn, c.CacheRedis),
        UserCouponsModel:   model.NewUserCouponsModel(conn, c.CacheRedis),
        PromotionsModel:    model.NewPromotionsModel(conn, c.CacheRedis),
        UserPointsModel:    model.NewUserPointsModel(conn, c.CacheRedis),
        PointsRecordsModel: model.NewPointsRecordsModel(conn, c.CacheRedis),
        Producer:           producer.NewMarketingProducer(ch, c.RabbitMQ.Exchange),
        MessageRpc:         messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc)),
    }
}