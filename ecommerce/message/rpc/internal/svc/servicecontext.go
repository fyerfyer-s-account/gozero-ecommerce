package svc

import "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
