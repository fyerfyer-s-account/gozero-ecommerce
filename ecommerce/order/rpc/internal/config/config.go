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

    OrderTimeout    int64 `json:",default=1800"`   // 30 minutes
    AutoConfirmTime int64 `json:",default=604800"` // 7 days
    DefaultPageSize int   `json:",default=10"`
    MaxOrderItems   int   `json:",default=50"`

    RabbitMQ struct {
        Host      string
        Port      int
        Username  string
        Password  string
        VHost     string
        Exchanges struct {
            OrderEvent struct {
                Name    string
                Type    string
                Durable bool
            }
        }
        Queues struct {
            OrderStatus struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            ShippingStatus struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
        }
    }

    UserRpc     zrpc.RpcClientConf
    CartRpc     zrpc.RpcClientConf
    ProductRpc  zrpc.RpcClientConf
}