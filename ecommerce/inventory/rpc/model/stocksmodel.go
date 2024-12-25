package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StocksModel = (*customStocksModel)(nil)

type (
	// StocksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStocksModel.
	StocksModel interface {
		stocksModel
		withSession(session sqlx.Session) StocksModel
	}

	customStocksModel struct {
		*defaultStocksModel
	}
)

// NewStocksModel returns a model for the database table.
func NewStocksModel(conn sqlx.SqlConn) StocksModel {
	return &customStocksModel{
		defaultStocksModel: newStocksModel(conn),
	}
}

func (m *customStocksModel) withSession(session sqlx.Session) StocksModel {
	return NewStocksModel(sqlx.NewSqlConnFromSession(session))
}
