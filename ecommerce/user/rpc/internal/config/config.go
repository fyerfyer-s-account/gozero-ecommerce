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

	Salt string

	JwtAuth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshRedis  struct {
			Host      string
			Type      string
			Pass      string
			KeyPrefix string
		}
	}

	PayTokenExpire int64
	PayTokenSecret string

	MinPasswordLength int
	MaxAddressCount   int
	InitialPoints     int64

	RabbitMQ struct {
		Host      string
		Port      int
		Username  string
		Password  string
		VHost     string
		Exchanges struct {
			UserEvent struct {
				Name    string
				Type    string
				Durable bool
			}
		}
		Queues struct {
			UserNotification struct {
				Name       string
				RoutingKey string
				Durable    bool
			}
		}
	}
}
