package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret   string
		AccessExpire   int64
		RefreshSecret  string
		RefreshExpire  int64
		BlacklistRedis struct {
			Host string
			Type string
			Pass string
			Key  string
		}
	}

	PageSize int

	AdminAuth struct {
		AccessSecret string
		AccessExpire int64
		RoleKey      string
	}

	Redis struct {
		Host string
		Type string
		Pass string
		Key  string
	}

	UserRpc    zrpc.RpcClientConf
	ProductRpc zrpc.RpcClientConf
	PaymentRpc zrpc.RpcClientConf
	CartRpc    zrpc.RpcClientConf
}
