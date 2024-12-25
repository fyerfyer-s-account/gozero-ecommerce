package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PaymentLogsModel = (*customPaymentLogsModel)(nil)

type (
	// PaymentLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentLogsModel.
	PaymentLogsModel interface {
		paymentLogsModel
		withSession(session sqlx.Session) PaymentLogsModel
	}

	customPaymentLogsModel struct {
		*defaultPaymentLogsModel
	}
)

// NewPaymentLogsModel returns a model for the database table.
func NewPaymentLogsModel(conn sqlx.SqlConn) PaymentLogsModel {
	return &customPaymentLogsModel{
		defaultPaymentLogsModel: newPaymentLogsModel(conn),
	}
}

func (m *customPaymentLogsModel) withSession(session sqlx.Session) PaymentLogsModel {
	return NewPaymentLogsModel(sqlx.NewSqlConnFromSession(session))
}
