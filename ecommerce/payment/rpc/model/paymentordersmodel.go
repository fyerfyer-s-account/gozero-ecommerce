package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentOrdersModel = (*customPaymentOrdersModel)(nil)

type (
	// PaymentOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentOrdersModel.
	PaymentOrdersModel interface {
		paymentOrdersModel
	}

	customPaymentOrdersModel struct {
		*defaultPaymentOrdersModel
	}
)

// NewPaymentOrdersModel returns a model for the database table.
func NewPaymentOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentOrdersModel {
	return &customPaymentOrdersModel{
		defaultPaymentOrdersModel: newPaymentOrdersModel(conn, c, opts...),
	}
}
