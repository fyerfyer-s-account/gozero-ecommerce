package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
	RabbitMQ   struct {
		Host      string
		Port      int
		Username  string
		Password  string
		VHost     string
		Exchanges map[string]struct {
			Name    string
			Type    string
			Durable bool
		}
		Queues map[string]struct {
			Name       string
			RoutingKey string
			Durable    bool
		}
	}
	
	ProductRpc         zrpc.RpcClientConf
	MaxItemsPerCart    int64
	MaxQuantityPerItem int64
	PageSize           int64
}
