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
            MessageEvent struct {
                Name    string
                Type    string
                Durable bool
            }
        }
        Queues struct {
            NotificationQueue struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            TemplateQueue struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
        }
    }
}