package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CartStatisticsModel = (*customCartStatisticsModel)(nil)

type (
	// CartStatisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartStatisticsModel.
	CartStatisticsModel interface {
		cartStatisticsModel
		withSession(session sqlx.Session) CartStatisticsModel
	}

	customCartStatisticsModel struct {
		*defaultCartStatisticsModel
	}
)

// NewCartStatisticsModel returns a model for the database table.
func NewCartStatisticsModel(conn sqlx.SqlConn) CartStatisticsModel {
	return &customCartStatisticsModel{
		defaultCartStatisticsModel: newCartStatisticsModel(conn),
	}
}

func (m *customCartStatisticsModel) withSession(session sqlx.Session) CartStatisticsModel {
	return NewCartStatisticsModel(sqlx.NewSqlConnFromSession(session))
}
