package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WarehousesModel = (*customWarehousesModel)(nil)

type (
	// WarehousesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWarehousesModel.
	WarehousesModel interface {
		warehousesModel
		withSession(session sqlx.Session) WarehousesModel
	}

	customWarehousesModel struct {
		*defaultWarehousesModel
	}
)

// NewWarehousesModel returns a model for the database table.
func NewWarehousesModel(conn sqlx.SqlConn) WarehousesModel {
	return &customWarehousesModel{
		defaultWarehousesModel: newWarehousesModel(conn),
	}
}

func (m *customWarehousesModel) withSession(session sqlx.Session) WarehousesModel {
	return NewWarehousesModel(sqlx.NewSqlConnFromSession(session))
}
