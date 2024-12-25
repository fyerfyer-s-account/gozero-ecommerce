package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PaymentChannelsModel = (*customPaymentChannelsModel)(nil)

type (
	// PaymentChannelsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentChannelsModel.
	PaymentChannelsModel interface {
		paymentChannelsModel
		withSession(session sqlx.Session) PaymentChannelsModel
	}

	customPaymentChannelsModel struct {
		*defaultPaymentChannelsModel
	}
)

// NewPaymentChannelsModel returns a model for the database table.
func NewPaymentChannelsModel(conn sqlx.SqlConn) PaymentChannelsModel {
	return &customPaymentChannelsModel{
		defaultPaymentChannelsModel: newPaymentChannelsModel(conn),
	}
}

func (m *customPaymentChannelsModel) withSession(session sqlx.Session) PaymentChannelsModel {
	return NewPaymentChannelsModel(sqlx.NewSqlConnFromSession(session))
}
