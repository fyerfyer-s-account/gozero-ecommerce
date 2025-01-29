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

    PaymentTimeout  int64   `json:",default=7200"`   
    RefundTimeout   int64   `json:",default=604800"` 
    MaxRetries      int     `json:",default=3"`
    DefaultPageSize int     `json:",default=10"`
    MaxAmount       float64 `json:",default=100000"`
    MinAmount       float64 `json:",default=0.01"`

    RabbitMQ struct {
        Host              string
        Port              int
        Username          string
        Password          string
        VHost             string
        Exchange          string
        ConnectionTimeout int
        HeartbeatInterval int
    }

    OrderRpc zrpc.RpcClientConf
    UserRpc  zrpc.RpcClientConf
}