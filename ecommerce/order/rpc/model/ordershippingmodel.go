package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderShippingModel = (*customOrderShippingModel)(nil)

type (
	// OrderShippingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderShippingModel.
	OrderShippingModel interface {
		orderShippingModel
	}

	customOrderShippingModel struct {
		*defaultOrderShippingModel
	}
)

// NewOrderShippingModel returns a model for the database table.
func NewOrderShippingModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderShippingModel {
	return &customOrderShippingModel{
		defaultOrderShippingModel: newOrderShippingModel(conn, c, opts...),
	}
}
