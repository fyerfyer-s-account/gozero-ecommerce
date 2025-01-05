package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AdminsModel = (*customAdminsModel)(nil)

type (
	// AdminsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminsModel.
	AdminsModel interface {
		adminsModel
		withSession(session sqlx.Session) AdminsModel
	}

	customAdminsModel struct {
		*defaultAdminsModel
	}
)

// NewAdminsModel returns a model for the database table.
func NewAdminsModel(conn sqlx.SqlConn) AdminsModel {
	return &customAdminsModel{
		defaultAdminsModel: newAdminsModel(conn),
	}
}

func (m *customAdminsModel) withSession(session sqlx.Session) AdminsModel {
	return NewAdminsModel(sqlx.NewSqlConnFromSession(session))
}
