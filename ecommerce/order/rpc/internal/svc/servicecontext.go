package svc

import (
    "context"
    "fmt"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cartclient"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/userclient"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config config.Config

    OrdersModel        model.OrdersModel
    OrderItemsModel    model.OrderItemsModel
    OrderShippingModel model.OrderShippingModel
    OrderRefundsModel  model.OrderRefundsModel
    OrderPaymentsModel model.OrderPaymentsModel

    UserRpc    userclient.User
    CartRpc    cartclient.Cart
    ProductRpc productservice.ProductService

    Producer *producer.OrderProducer
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

    // Initialize models and services
    return &ServiceContext{
        Config: c,

        OrdersModel:        model.NewOrdersModel(conn, c.CacheRedis),
        OrderItemsModel:    model.NewOrderItemsModel(conn, c.CacheRedis),
        OrderShippingModel: model.NewOrderShippingModel(conn, c.CacheRedis),
        OrderRefundsModel:  model.NewOrderRefundsModel(conn, c.CacheRedis),
        OrderPaymentsModel: model.NewOrderPaymentsModel(conn, c.CacheRedis),

        UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
        CartRpc:    cartclient.NewCart(zrpc.MustNewClient(c.CartRpc)),
        ProductRpc: productservice.NewProductService(zrpc.MustNewClient(c.ProductRpc)),
        Producer:   producer.NewOrderProducer(ch, c.RabbitMQ.Exchange),
    }
}