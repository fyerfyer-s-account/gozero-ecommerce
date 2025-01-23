package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	RabbitMQ          struct {
        Host      string
        Port      int
        Username  string
        Password  string
        VHost     string
        Exchanges struct {
            InventoryEvent struct {
                Name    string
                Type    string
                Durable bool
            }
        }
        Queues struct {
            StockUpdate struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            StockAlert struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            StockLock struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
            OrderEvents struct {
                Name       string
                RoutingKey string
                Durable    bool
            }
        }
        Retry struct {
            MaxAttempts     int
            InitialInterval int
            MaxInterval     int
            BackoffFactor   float64
            Jitter         bool
        }
        Batch struct {
            Size          int
            FlushInterval int
            Workers       int
        }
        DeadLetter struct {
            Exchange   string
            Queue      string
            RoutingKey string
        }
        Middleware struct {
            EnableRecovery bool
            EnableLogging  bool
        }
        Server struct {
            Name      string
            Mode      string
            LogLevel  string
            Consumers []struct {
                Queue   string
                Workers int
            }
            Monitor struct {
                Enabled bool
                Port    int
            }
        }
    }
	MessageRpc   zrpc.RpcClientConf
	InventoryRpc zrpc.RpcClientConf
}
