package svc

import "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
