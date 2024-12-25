package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SkusModel = (*customSkusModel)(nil)

type (
	// SkusModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSkusModel.
	SkusModel interface {
		skusModel
		withSession(session sqlx.Session) SkusModel
	}

	customSkusModel struct {
		*defaultSkusModel
	}
)

// NewSkusModel returns a model for the database table.
func NewSkusModel(conn sqlx.SqlConn) SkusModel {
	return &customSkusModel{
		defaultSkusModel: newSkusModel(conn),
	}
}

func (m *customSkusModel) withSession(session sqlx.Session) SkusModel {
	return NewSkusModel(sqlx.NewSqlConnFromSession(session))
}
