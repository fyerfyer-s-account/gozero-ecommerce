package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	RabbitMQ     RabbitMQConfig `yaml:"RabbitMQ"`
	MessageRpc   zrpc.RpcClientConf
	InventoryRpc zrpc.RpcClientConf
}
