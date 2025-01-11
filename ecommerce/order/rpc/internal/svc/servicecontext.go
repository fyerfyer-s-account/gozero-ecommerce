package svc

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cartclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config config.Config

    OrdersModel         model.OrdersModel
    OrderItemsModel     model.OrderItemsModel
    OrderShippingModel  model.OrderShippingModel
    OrderRefundsModel   model.OrderRefundsModel
    OrderPaymentsModel  model.OrderPaymentsModel
    
    UserRpc            userclient.User
    CartRpc            cartclient.Cart
    ProductRpc         productservice.ProductService

    // OrderEventProducer *rabbitmq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)  
    // orderEventProducer := rabbitmq.New.NewProducer(
    //     []string{c.RabbitMQ.Host}, 
    //     c.RabbitMQ.Username, 
    //     c.RabbitMQ.Password, 
    //     c.RabbitMQ.VHost,
    // )

    return &ServiceContext{
        Config: c,

        OrdersModel:         model.NewOrdersModel(conn, c.CacheRedis),
        OrderItemsModel:     model.NewOrderItemsModel(conn, c.CacheRedis),
        OrderShippingModel:  model.NewOrderShippingModel(conn, c.CacheRedis),
        OrderRefundsModel:   model.NewOrderRefundsModel(conn, c.CacheRedis),
        OrderPaymentsModel:  model.NewOrderPaymentsModel(conn, c.CacheRedis),

        UserRpc:            userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
        CartRpc:            cartclient.NewCart(zrpc.MustNewClient(c.CartRpc)),
        ProductRpc:         productservice.NewProductService(zrpc.MustNewClient(c.ProductRpc)),
        // OrderEventProducer: orderEventProducer,
    }
}