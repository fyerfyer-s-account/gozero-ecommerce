package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
	}

	Redis struct {
		Host string
		Type string
		Pass string
		Key  string
	}

	UserRpc      zrpc.RpcClientConf
	ProductRpc   zrpc.RpcClientConf
	CartRpc      zrpc.RpcClientConf
	OrderRpc     zrpc.RpcClientConf
	PaymentRpc   zrpc.RpcClientConf
	InventoryRpc zrpc.RpcClientConf
	MarketingRpc zrpc.RpcClientConf
	SearchRpc    zrpc.RpcClientConf
	MessageRpc   zrpc.RpcClientConf
}
