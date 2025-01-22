package config

import (
	rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	CacheRedis         cache.CacheConf

	RabbitMQ           rmqconfig.RabbitMQConfig
	
	ProductRpc         zrpc.RpcClientConf
	MaxItemsPerCart    int64
	MaxQuantityPerItem int64
	PageSize           int64
}
