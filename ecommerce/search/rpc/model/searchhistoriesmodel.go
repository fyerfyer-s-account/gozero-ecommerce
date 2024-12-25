package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SearchHistoriesModel = (*customSearchHistoriesModel)(nil)

type (
	// SearchHistoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchHistoriesModel.
	SearchHistoriesModel interface {
		searchHistoriesModel
		withSession(session sqlx.Session) SearchHistoriesModel
	}

	customSearchHistoriesModel struct {
		*defaultSearchHistoriesModel
	}
)

// NewSearchHistoriesModel returns a model for the database table.
func NewSearchHistoriesModel(conn sqlx.SqlConn) SearchHistoriesModel {
	return &customSearchHistoriesModel{
		defaultSearchHistoriesModel: newSearchHistoriesModel(conn),
	}
}

func (m *customSearchHistoriesModel) withSession(session sqlx.Session) SearchHistoriesModel {
	return NewSearchHistoriesModel(sqlx.NewSqlConnFromSession(session))
}
