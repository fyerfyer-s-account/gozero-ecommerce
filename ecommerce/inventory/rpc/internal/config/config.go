package config

import (
	rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis         cache.CacheConf
	StockLockTimeout   int
	MaxBatchSize       int
	AlertThreshold     int
	DefaultWarehouseId int64
	PageSize           int
	RabbitMQ           rmqconfig.RabbitMQConfig
	MessageRpc         zrpc.RpcClientConf
	Etcd               struct {
		Hosts []string
		Key   string
		zrpc.RpcClientConf
	}
}
