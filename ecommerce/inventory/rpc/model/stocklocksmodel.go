package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StockLocksModel = (*customStockLocksModel)(nil)

type (
	// StockLocksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockLocksModel.
	StockLocksModel interface {
		stockLocksModel
		withSession(session sqlx.Session) StockLocksModel
	}

	customStockLocksModel struct {
		*defaultStockLocksModel
	}
)

// NewStockLocksModel returns a model for the database table.
func NewStockLocksModel(conn sqlx.SqlConn) StockLocksModel {
	return &customStockLocksModel{
		defaultStockLocksModel: newStockLocksModel(conn),
	}
}

func (m *customStockLocksModel) withSession(session sqlx.Session) StockLocksModel {
	return NewStockLocksModel(sqlx.NewSqlConnFromSession(session))
}
