package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefundOrdersModel = (*customRefundOrdersModel)(nil)

type (
	// RefundOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundOrdersModel.
	RefundOrdersModel interface {
		refundOrdersModel
	}

	customRefundOrdersModel struct {
		*defaultRefundOrdersModel
	}
)

// NewRefundOrdersModel returns a model for the database table.
func NewRefundOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RefundOrdersModel {
	return &customRefundOrdersModel{
		defaultRefundOrdersModel: newRefundOrdersModel(conn, c, opts...),
	}
}
