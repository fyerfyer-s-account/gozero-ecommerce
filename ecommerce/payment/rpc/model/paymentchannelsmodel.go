package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentChannelsModel = (*customPaymentChannelsModel)(nil)

type (
	// PaymentChannelsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentChannelsModel.
	PaymentChannelsModel interface {
		paymentChannelsModel
	}

	customPaymentChannelsModel struct {
		*defaultPaymentChannelsModel
	}
)

// NewPaymentChannelsModel returns a model for the database table.
func NewPaymentChannelsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentChannelsModel {
	return &customPaymentChannelsModel{
		defaultPaymentChannelsModel: newPaymentChannelsModel(conn, c, opts...),
	}
}
