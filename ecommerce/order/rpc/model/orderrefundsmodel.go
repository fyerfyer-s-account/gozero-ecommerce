package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderRefundsModel = (*customOrderRefundsModel)(nil)

type (
	// OrderRefundsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderRefundsModel.
	OrderRefundsModel interface {
		orderRefundsModel
	}

	customOrderRefundsModel struct {
		*defaultOrderRefundsModel
	}
)

// NewOrderRefundsModel returns a model for the database table.
func NewOrderRefundsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderRefundsModel {
	return &customOrderRefundsModel{
		defaultOrderRefundsModel: newOrderRefundsModel(conn, c, opts...),
	}
}
