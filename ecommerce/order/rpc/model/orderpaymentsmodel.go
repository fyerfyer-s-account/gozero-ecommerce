package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderPaymentsModel = (*customOrderPaymentsModel)(nil)

type (
	// OrderPaymentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderPaymentsModel.
	OrderPaymentsModel interface {
		orderPaymentsModel
	}

	customOrderPaymentsModel struct {
		*defaultOrderPaymentsModel
	}
)

// NewOrderPaymentsModel returns a model for the database table.
func NewOrderPaymentsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderPaymentsModel {
	return &customOrderPaymentsModel{
		defaultOrderPaymentsModel: newOrderPaymentsModel(conn, c, opts...),
	}
}
