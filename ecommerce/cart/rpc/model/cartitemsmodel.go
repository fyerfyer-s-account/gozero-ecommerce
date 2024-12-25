package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CartItemsModel = (*customCartItemsModel)(nil)

type (
	// CartItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartItemsModel.
	CartItemsModel interface {
		cartItemsModel
		withSession(session sqlx.Session) CartItemsModel
	}

	customCartItemsModel struct {
		*defaultCartItemsModel
	}
)

// NewCartItemsModel returns a model for the database table.
func NewCartItemsModel(conn sqlx.SqlConn) CartItemsModel {
	return &customCartItemsModel{
		defaultCartItemsModel: newCartItemsModel(conn),
	}
}

func (m *customCartItemsModel) withSession(session sqlx.Session) CartItemsModel {
	return NewCartItemsModel(sqlx.NewSqlConnFromSession(session))
}
