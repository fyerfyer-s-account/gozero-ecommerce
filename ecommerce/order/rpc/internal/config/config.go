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

    OrderTimeout    int64
    AutoConfirmTime int64
    DefaultPageSize int
    MaxOrderItems   int

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

    UserRpc    zrpc.RpcClientConf
    CartRpc    zrpc.RpcClientConf
    ProductRpc zrpc.RpcClientConf
}