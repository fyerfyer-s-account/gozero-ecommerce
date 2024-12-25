package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StockRecordsModel = (*customStockRecordsModel)(nil)

type (
	// StockRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockRecordsModel.
	StockRecordsModel interface {
		stockRecordsModel
		withSession(session sqlx.Session) StockRecordsModel
	}

	customStockRecordsModel struct {
		*defaultStockRecordsModel
	}
)

// NewStockRecordsModel returns a model for the database table.
func NewStockRecordsModel(conn sqlx.SqlConn) StockRecordsModel {
	return &customStockRecordsModel{
		defaultStockRecordsModel: newStockRecordsModel(conn),
	}
}

func (m *customStockRecordsModel) withSession(session sqlx.Session) StockRecordsModel {
	return NewStockRecordsModel(sqlx.NewSqlConnFromSession(session))
}
