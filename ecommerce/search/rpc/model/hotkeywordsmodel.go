package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HotKeywordsModel = (*customHotKeywordsModel)(nil)

type (
	// HotKeywordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHotKeywordsModel.
	HotKeywordsModel interface {
		hotKeywordsModel
		withSession(session sqlx.Session) HotKeywordsModel
	}

	customHotKeywordsModel struct {
		*defaultHotKeywordsModel
	}
)

// NewHotKeywordsModel returns a model for the database table.
func NewHotKeywordsModel(conn sqlx.SqlConn) HotKeywordsModel {
	return &customHotKeywordsModel{
		defaultHotKeywordsModel: newHotKeywordsModel(conn),
	}
}

func (m *customHotKeywordsModel) withSession(session sqlx.Session) HotKeywordsModel {
	return NewHotKeywordsModel(sqlx.NewSqlConnFromSession(session))
}
