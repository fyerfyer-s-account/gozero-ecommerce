package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ProductIndicesModel = (*customProductIndicesModel)(nil)

type (
	// ProductIndicesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductIndicesModel.
	ProductIndicesModel interface {
		productIndicesModel
		withSession(session sqlx.Session) ProductIndicesModel
	}

	customProductIndicesModel struct {
		*defaultProductIndicesModel
	}
)

// NewProductIndicesModel returns a model for the database table.
func NewProductIndicesModel(conn sqlx.SqlConn) ProductIndicesModel {
	return &customProductIndicesModel{
		defaultProductIndicesModel: newProductIndicesModel(conn),
	}
}

func (m *customProductIndicesModel) withSession(session sqlx.Session) ProductIndicesModel {
	return NewProductIndicesModel(sqlx.NewSqlConnFromSession(session))
}
