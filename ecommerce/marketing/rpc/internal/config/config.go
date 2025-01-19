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
            MarketingEvent struct {
                Name    string
                Type    string
                Durable bool
            }
        }
        Queues struct {
            CouponEvent struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            PromotionEvent struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            PointsEvent struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
        }
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