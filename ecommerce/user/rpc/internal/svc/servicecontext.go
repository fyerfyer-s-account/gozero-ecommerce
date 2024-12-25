package svc

import "github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
