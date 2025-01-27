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

    // Inventory specific settings
    StockLockTimeout   int64
    MaxBatchSize       int
    AlertThreshold     int32
    DefaultWarehouseId int64
    DefaultPageSize    int

    // RabbitMQ Configuration
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

    MessageRpc zrpc.RpcClientConf
}