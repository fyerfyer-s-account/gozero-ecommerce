package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentLogsModel = (*customPaymentLogsModel)(nil)

type (
	// PaymentLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentLogsModel.
	PaymentLogsModel interface {
		paymentLogsModel
	}

	customPaymentLogsModel struct {
		*defaultPaymentLogsModel
	}
)

// NewPaymentLogsModel returns a model for the database table.
func NewPaymentLogsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentLogsModel {
	return &customPaymentLogsModel{
		defaultPaymentLogsModel: newPaymentLogsModel(conn, c, opts...),
	}
}
