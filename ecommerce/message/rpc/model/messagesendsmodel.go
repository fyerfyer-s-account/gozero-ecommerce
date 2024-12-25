package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ MessageSendsModel = (*customMessageSendsModel)(nil)

type (
	// MessageSendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageSendsModel.
	MessageSendsModel interface {
		messageSendsModel
		withSession(session sqlx.Session) MessageSendsModel
	}

	customMessageSendsModel struct {
		*defaultMessageSendsModel
	}
)

// NewMessageSendsModel returns a model for the database table.
func NewMessageSendsModel(conn sqlx.SqlConn) MessageSendsModel {
	return &customMessageSendsModel{
		defaultMessageSendsModel: newMessageSendsModel(conn),
	}
}

func (m *customMessageSendsModel) withSession(session sqlx.Session) MessageSendsModel {
	return NewMessageSendsModel(sqlx.NewSqlConnFromSession(session))
}
