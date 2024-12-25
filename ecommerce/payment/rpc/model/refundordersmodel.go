package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RefundOrdersModel = (*customRefundOrdersModel)(nil)

type (
	// RefundOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundOrdersModel.
	RefundOrdersModel interface {
		refundOrdersModel
		withSession(session sqlx.Session) RefundOrdersModel
	}

	customRefundOrdersModel struct {
		*defaultRefundOrdersModel
	}
)

// NewRefundOrdersModel returns a model for the database table.
func NewRefundOrdersModel(conn sqlx.SqlConn) RefundOrdersModel {
	return &customRefundOrdersModel{
		defaultRefundOrdersModel: newRefundOrdersModel(conn),
	}
}

func (m *customRefundOrdersModel) withSession(session sqlx.Session) RefundOrdersModel {
	return NewRefundOrdersModel(sqlx.NewSqlConnFromSession(session))
}
