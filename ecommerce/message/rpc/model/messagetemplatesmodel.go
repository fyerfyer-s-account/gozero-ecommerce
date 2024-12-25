package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ MessageTemplatesModel = (*customMessageTemplatesModel)(nil)

type (
	// MessageTemplatesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageTemplatesModel.
	MessageTemplatesModel interface {
		messageTemplatesModel
		withSession(session sqlx.Session) MessageTemplatesModel
	}

	customMessageTemplatesModel struct {
		*defaultMessageTemplatesModel
	}
)

// NewMessageTemplatesModel returns a model for the database table.
func NewMessageTemplatesModel(conn sqlx.SqlConn) MessageTemplatesModel {
	return &customMessageTemplatesModel{
		defaultMessageTemplatesModel: newMessageTemplatesModel(conn),
	}
}

func (m *customMessageTemplatesModel) withSession(session sqlx.Session) MessageTemplatesModel {
	return NewMessageTemplatesModel(sqlx.NewSqlConnFromSession(session))
}
