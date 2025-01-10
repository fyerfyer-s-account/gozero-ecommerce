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

	PaymentTimeout  int64   `json:",default=7200"`   // 2 hours
	RefundTimeout   int64   `json:",default=604800"` // 7 days
	MaxRetries      int     `json:",default=3"`
	DefaultPageSize int     `json:",default=10"`
	MaxAmount       float64 `json:",default=100000"`
	MinAmount       float64 `json:",default=0.01"`

	RabbitMQ struct {
		Host      string
		Port      int
		Username  string
		Password  string
		VHost     string
		Exchanges struct {
			PaymentEvent struct {
				Name    string
				Type    string
				Durable bool
			}
		}
		Queues struct {
			PaymentStatus struct {
				Name       string
				RoutingKey string
				Durable    bool
			}
			RefundStatus struct {
				Name       string
				RoutingKey string
				Durable    bool
			}
		}
	}

	OrderRpc   zrpc.RpcClientConf
	UserRpc zrpc.RpcClientConf
}
