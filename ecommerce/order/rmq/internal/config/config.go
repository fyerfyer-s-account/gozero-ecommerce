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
		Host              string `json:"host"`
		Port              int    `json:"port"`
		Username          string `json:"username"`
		Password          string `json:"password"`
		VHost             string `json:"vhost"`
		ConnectionTimeout int    `json:"connection_timeout"`
		HeartbeatInterval int    `json:"heartbeat_interval"`
		PrefetchCount     int    `json:"prefetch_count"`
		PrefetchGlobal    bool   `json:"prefetch_global"`
		Exchanges         []struct {
			Name       string `json:"name"`
			Type       string `json:"type"`
			Durable    bool   `json:"durable"`
			AutoDelete bool   `json:"auto_delete"`
			Internal   bool   `json:"internal"`
			NoWait     bool   `json:"no_wait"`
		} `json:"exchanges"`
		Queues []struct {
			Name       string `json:"name"`
			Durable    bool   `json:"durable"`
			AutoDelete bool   `json:"auto_delete"`
			Exclusive  bool   `json:"exclusive"`
			NoWait     bool   `json:"no_wait"`
			Bindings   []struct {
				Exchange   string `json:"exchange"`
				RoutingKey string `json:"routing_key"`
				NoWait     bool   `json:"no_wait"`
			} `json:"bindings"`
		} `json:"queues"`
	}
}
