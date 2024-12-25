package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LoginRecordsModel = (*customLoginRecordsModel)(nil)

type (
	// LoginRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginRecordsModel.
	LoginRecordsModel interface {
		loginRecordsModel
		withSession(session sqlx.Session) LoginRecordsModel
	}

	customLoginRecordsModel struct {
		*defaultLoginRecordsModel
	}
)

// NewLoginRecordsModel returns a model for the database table.
func NewLoginRecordsModel(conn sqlx.SqlConn) LoginRecordsModel {
	return &customLoginRecordsModel{
		defaultLoginRecordsModel: newLoginRecordsModel(conn),
	}
}

func (m *customLoginRecordsModel) withSession(session sqlx.Session) LoginRecordsModel {
	return NewLoginRecordsModel(sqlx.NewSqlConnFromSession(session))
}
