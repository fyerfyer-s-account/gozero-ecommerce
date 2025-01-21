package config

import (
	rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	OrderTimeout    int64
	AutoConfirmTime int64
	DefaultPageSize int
	MaxOrderItems   int

	RabbitMQ rmqconfig.RabbitMQConfig

	UserRpc    zrpc.RpcClientConf
	CartRpc    zrpc.RpcClientConf
	ProductRpc zrpc.RpcClientConf
}
