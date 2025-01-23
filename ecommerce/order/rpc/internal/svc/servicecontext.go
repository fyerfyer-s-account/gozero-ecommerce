package svc

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cartclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/logx"
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

	Producer *producer.Producer
	Consumer *consumer.Consumer // Add consumer field for RMQ server
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	logger := &zerolog.LogWrapper{Logger: logx.WithContext(context.TODO())}

	prod, err := producer.NewProducer(&c.RabbitMQ)
	if err != nil {
		panic(err)
	}

	// Initialize consumer if in RMQ server mode
	var cons *consumer.Consumer
	if c.RabbitMQ.Server.Mode != "" {
		cons, err = consumer.NewConsumer(&c.RabbitMQ, logger)
		if err != nil {
			panic(err)
		}
	}

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
		Producer:   prod,
		Consumer:   cons,
	}
}
