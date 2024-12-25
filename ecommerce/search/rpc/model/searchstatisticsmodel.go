package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SearchStatisticsModel = (*customSearchStatisticsModel)(nil)

type (
	// SearchStatisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSearchStatisticsModel.
	SearchStatisticsModel interface {
		searchStatisticsModel
		withSession(session sqlx.Session) SearchStatisticsModel
	}

	customSearchStatisticsModel struct {
		*defaultSearchStatisticsModel
	}
)

// NewSearchStatisticsModel returns a model for the database table.
func NewSearchStatisticsModel(conn sqlx.SqlConn) SearchStatisticsModel {
	return &customSearchStatisticsModel{
		defaultSearchStatisticsModel: newSearchStatisticsModel(conn),
	}
}

func (m *customSearchStatisticsModel) withSession(session sqlx.Session) SearchStatisticsModel {
	return NewSearchStatisticsModel(sqlx.NewSqlConnFromSession(session))
}
