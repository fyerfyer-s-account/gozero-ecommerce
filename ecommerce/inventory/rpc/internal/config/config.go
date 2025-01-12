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
		Exchanges struct {
			InventoryEvent struct {
				Name    string
				Type    string
				Durable bool
			}
		}
		Queues struct {
			StockUpdate struct {
				Name       string
				RoutingKey string
				Durable    bool
			}
			StockAlert struct {
				Name       string
				RoutingKey string
				Durable    bool
			}
		}
	}
	StockLockTimeout int   `json:",default=1800"`
	MaxBatchSize     int   `json:",default=100"`
	AlertThreshold   int   `json:",default=10"`
	DefaultWarehouse int64 `json:",default=1"`
	PageSize         int   `json:",default=20"`
}
