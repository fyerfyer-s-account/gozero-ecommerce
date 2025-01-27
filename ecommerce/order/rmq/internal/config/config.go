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

	RabbitMQ struct {
		Host              string
		Port              int
		Username          string
		Password          string
		VHost             string
		ConnectionTimeout int
		HeartbeatInterval int
		PrefetchCount     int
		PrefetchGlobal    bool
		Exchanges         []struct {
			Name       string
			Type       string
			Durable    bool
			AutoDelete bool
			Internal   bool
			NoWait     bool
		}
		Queues []struct {
			Name       string
			Durable    bool
			AutoDelete bool
			Exclusive  bool
			NoWait     bool
			Bindings   []struct {
				Exchange   string
				RoutingKey string
				NoWait     bool
			}
		}
	}
}
