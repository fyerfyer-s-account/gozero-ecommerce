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
        Exchange          string
        ConnectionTimeout int
        HeartbeatInterval int
    }

    CouponLimits struct {
        MaxPerUser  int
        MaxActive   int
        BatchSize   int
    }

    PromotionLimits struct {
        MaxActive int
        MaxRules  int
    }

    PointsLimits struct {
        MaxPoints   int64
        MinPoints   int64
        ExpireDays  int
    }

    MessageRpc zrpc.RpcClientConf
}