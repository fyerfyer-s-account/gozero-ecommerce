package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PointsRecordsModel = (*customPointsRecordsModel)(nil)

type (
	// PointsRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPointsRecordsModel.
	PointsRecordsModel interface {
		pointsRecordsModel
		withSession(session sqlx.Session) PointsRecordsModel
	}

	customPointsRecordsModel struct {
		*defaultPointsRecordsModel
	}
)

// NewPointsRecordsModel returns a model for the database table.
func NewPointsRecordsModel(conn sqlx.SqlConn) PointsRecordsModel {
	return &customPointsRecordsModel{
		defaultPointsRecordsModel: newPointsRecordsModel(conn),
	}
}

func (m *customPointsRecordsModel) withSession(session sqlx.Session) PointsRecordsModel {
	return NewPointsRecordsModel(sqlx.NewSqlConnFromSession(session))
}
