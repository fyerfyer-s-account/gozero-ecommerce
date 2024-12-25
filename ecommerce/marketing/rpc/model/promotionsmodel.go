package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PromotionsModel = (*customPromotionsModel)(nil)

type (
	// PromotionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPromotionsModel.
	PromotionsModel interface {
		promotionsModel
		withSession(session sqlx.Session) PromotionsModel
	}

	customPromotionsModel struct {
		*defaultPromotionsModel
	}
)

// NewPromotionsModel returns a model for the database table.
func NewPromotionsModel(conn sqlx.SqlConn) PromotionsModel {
	return &customPromotionsModel{
		defaultPromotionsModel: newPromotionsModel(conn),
	}
}

func (m *customPromotionsModel) withSession(session sqlx.Session) PromotionsModel {
	return NewPromotionsModel(sqlx.NewSqlConnFromSession(session))
}
