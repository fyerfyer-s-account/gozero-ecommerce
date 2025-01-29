package svc

import (
    "context"
    "fmt"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/orderservice"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/userclient"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config config.Config

    PaymentOrdersModel   model.PaymentOrdersModel
    RefundOrdersModel    model.RefundOrdersModel
    PaymentChannelsModel model.PaymentChannelsModel
    PaymentLogsModel    model.PaymentLogsModel

    OrderRpc orderservice.OrderService
    UserRpc  userclient.User

    Producer *producer.PaymentProducer
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
        Config:               c,
        PaymentOrdersModel:   model.NewPaymentOrdersModel(conn, c.CacheRedis),
        RefundOrdersModel:    model.NewRefundOrdersModel(conn, c.CacheRedis),
        PaymentChannelsModel: model.NewPaymentChannelsModel(conn, c.CacheRedis),
        PaymentLogsModel:     model.NewPaymentLogsModel(conn, c.CacheRedis),
        OrderRpc:             orderservice.NewOrderService(zrpc.MustNewClient(c.OrderRpc)),
        UserRpc:              userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
        Producer:            producer.NewPaymentProducer(ch, c.RabbitMQ.Exchange),
    }
}