package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WarehousesModel = (*customWarehousesModel)(nil)

type (
	// WarehousesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWarehousesModel.
	WarehousesModel interface {
		warehousesModel
		FindMany(ctx context.Context, page, pageSize int32) ([]*Warehouses, error)
        Count(ctx context.Context) (int64, error)
        UpdateStatus(ctx context.Context, id uint64, status int64) error
        FindByName(ctx context.Context, name string) (*Warehouses, error)
        FindManyByStatus(ctx context.Context, status int64) ([]*Warehouses, error)
	}

	customWarehousesModel struct {
		*defaultWarehousesModel
	}
)

// NewWarehousesModel returns a model for the database table.
func NewWarehousesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WarehousesModel {
	return &customWarehousesModel{
		defaultWarehousesModel: newWarehousesModel(conn, c, opts...),
	}
}

func (m *customWarehousesModel) FindMany(ctx context.Context, page, pageSize int32) ([]*Warehouses, error) {
    var warehouses []*Warehouses
    query := fmt.Sprintf("select %s from %s order by id desc limit ?, ?", warehousesRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &warehouses, query, (page-1)*pageSize, pageSize)
    return warehouses, err
}

func (m *customWarehousesModel) Count(ctx context.Context) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query)
    return count, err
}

func (m *customWarehousesModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
    mallInventoryWarehousesIdKey := fmt.Sprintf("%s%v", cacheMallInventoryWarehousesIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, id)
    }, mallInventoryWarehousesIdKey)
    return err
}

func (m *customWarehousesModel) FindByName(ctx context.Context, name string) (*Warehouses, error) {
    var warehouse Warehouses
    query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", warehousesRows, m.table)
    err := m.QueryRowNoCacheCtx(ctx, &warehouse, query, name)
    switch err {
    case nil:
        return &warehouse, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *customWarehousesModel) FindManyByStatus(ctx context.Context, status int64) ([]*Warehouses, error) {
    var warehouses []*Warehouses
    query := fmt.Sprintf("select %s from %s where `status` = ?", warehousesRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &warehouses, query, status)
    return warehouses, err
}