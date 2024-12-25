package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrderRefundsModel = (*customOrderRefundsModel)(nil)

type (
	// OrderRefundsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderRefundsModel.
	OrderRefundsModel interface {
		orderRefundsModel
		withSession(session sqlx.Session) OrderRefundsModel
	}

	customOrderRefundsModel struct {
		*defaultOrderRefundsModel
	}
)

// NewOrderRefundsModel returns a model for the database table.
func NewOrderRefundsModel(conn sqlx.SqlConn) OrderRefundsModel {
	return &customOrderRefundsModel{
		defaultOrderRefundsModel: newOrderRefundsModel(conn),
	}
}

func (m *customOrderRefundsModel) withSession(session sqlx.Session) OrderRefundsModel {
	return NewOrderRefundsModel(sqlx.NewSqlConnFromSession(session))
}
