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

    PaymentTimeout  int64      
    RefundTimeout   int64   
    MaxRetries      int     
    DefaultPageSize int     
    MaxAmount       float64 
    MinAmount       float64 

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