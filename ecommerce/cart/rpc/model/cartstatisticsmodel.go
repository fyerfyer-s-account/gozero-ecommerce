package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CartStatisticsModel = (*customCartStatisticsModel)(nil)

type (
	// CartStatisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartStatisticsModel.
	CartStatisticsModel interface {
		cartStatisticsModel
	}

	customCartStatisticsModel struct {
		*defaultCartStatisticsModel
	}
)

// NewCartStatisticsModel returns a model for the database table.
func NewCartStatisticsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CartStatisticsModel {
	return &customCartStatisticsModel{
		defaultCartStatisticsModel: newCartStatisticsModel(conn, c, opts...),
	}
}
