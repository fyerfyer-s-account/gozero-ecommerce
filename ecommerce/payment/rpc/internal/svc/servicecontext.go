package svc

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/orderservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	// Payment models
	PaymentOrdersModel   model.PaymentOrdersModel
	RefundOrdersModel    model.RefundOrdersModel
	PaymentChannelsModel model.PaymentChannelsModel
	PaymentLogsModel     model.PaymentLogsModel
	OrderRpc             orderservice.OrderService
	UserRpc              userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:               c,
		PaymentOrdersModel:   model.NewPaymentOrdersModel(conn, c.CacheRedis),
		RefundOrdersModel:    model.NewRefundOrdersModel(conn, c.CacheRedis),
		PaymentChannelsModel: model.NewPaymentChannelsModel(conn, c.CacheRedis),
		PaymentLogsModel:     model.NewPaymentLogsModel(conn, c.CacheRedis),
		OrderRpc:             orderservice.NewOrderService(zrpc.MustNewClient(c.OrderRpc)),
		UserRpc:              userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
