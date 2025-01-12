package model

import (
	"context"
    "database/sql"
    "fmt"
    "strings"
	
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StocksModel = (*customStocksModel)(nil)

type (
	// StocksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStocksModel.
	StocksModel interface {
		stocksModel
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
        BatchGet(ctx context.Context, skuIds []uint64, warehouseId uint64) (map[uint64]*Stocks, error)
        IncrAvailable(ctx context.Context, skuId, warehouseId uint64, quantity int64) error
        DecrAvailable(ctx context.Context, skuId, warehouseId uint64, quantity int64) error
        Lock(ctx context.Context, session sqlx.Session, skuId, warehouseId uint64, quantity int64) error
        Unlock(ctx context.Context, session sqlx.Session, skuId, warehouseId uint64, quantity int64) error
        Deduct(ctx context.Context, session sqlx.Session, skuId, warehouseId uint64, quantity int64) error
        FindManyByWarehouse(ctx context.Context, warehouseId uint64, page, pageSize int32) ([]*Stocks, error)
        Count(ctx context.Context, warehouseId uint64) (int64, error)
	}

	customStocksModel struct {
		*defaultStocksModel
	}
)

// NewStocksModel returns a model for the database table.
func NewStocksModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StocksModel {
	return &customStocksModel{
		defaultStocksModel: newStocksModel(conn, c, opts...),
	}
}

func (m *customStocksModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, fn)
}

// Batch get stocks
func (m *customStocksModel) BatchGet(ctx context.Context, skuIds []uint64, warehouseId uint64) (map[uint64]*Stocks, error) {
    if len(skuIds) == 0 {
        return make(map[uint64]*Stocks), nil
    }

    // Convert skuIds to string for IN clause
    ids := make([]string, len(skuIds))
    for i, id := range skuIds {
        ids[i] = fmt.Sprint(id)
    }

    query := fmt.Sprintf("select %s from %s where `sku_id` in (%s) and `warehouse_id` = ?",
        stocksRows, m.table, strings.Join(ids, ","))

    var stocks []*Stocks
    err := m.QueryRowsNoCacheCtx(ctx, &stocks, query, warehouseId)
    if err != nil {
        return nil, err
    }

    result := make(map[uint64]*Stocks)
    for _, stock := range stocks {
        result[stock.SkuId] = stock
    }
    return result, nil
}

// Increase available stock
func (m *customStocksModel) IncrAvailable(ctx context.Context, skuId, warehouseId uint64, quantity int64) error {
    key := fmt.Sprintf("%s%v:%v", cacheMallInventoryStocksSkuIdWarehouseIdPrefix, skuId, warehouseId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `available` = `available` + ?, `total` = `total` + ? where `sku_id` = ? and `warehouse_id` = ?", m.table)
        return conn.ExecCtx(ctx, query, quantity, quantity, skuId, warehouseId)
    }, key)
    return err
}

// Decrease available stock
func (m *customStocksModel) DecrAvailable(ctx context.Context, skuId, warehouseId uint64, quantity int64) error {
    key := fmt.Sprintf("%s%v:%v", cacheMallInventoryStocksSkuIdWarehouseIdPrefix, skuId, warehouseId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `available` = `available` - ?, `total` = `total` - ? where `sku_id` = ? and `warehouse_id` = ? and `available` >= ?", m.table)
        return conn.ExecCtx(ctx, query, quantity, quantity, skuId, warehouseId, quantity)
    }, key)
    return err
}

// Lock stock
func (m *customStocksModel) Lock(ctx context.Context, session sqlx.Session, skuId, warehouseId uint64, quantity int64) error {
    query := fmt.Sprintf("update %s set `available` = `available` - ?, `locked` = `locked` + ? where `sku_id` = ? and `warehouse_id` = ? and `available` >= ?", m.table)
    result, err := session.ExecCtx(ctx, query, quantity, quantity, skuId, warehouseId, quantity)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return ErrNotFound
    }
    return nil
}

// Unlock stock
func (m *customStocksModel) Unlock(ctx context.Context, session sqlx.Session, skuId, warehouseId uint64, quantity int64) error {
    query := fmt.Sprintf("update %s set `available` = `available` + ?, `locked` = `locked` - ? where `sku_id` = ? and `warehouse_id` = ? and `locked` >= ?", m.table)
    result, err := session.ExecCtx(ctx, query, quantity, quantity, skuId, warehouseId, quantity)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return ErrNotFound
    }
    return nil
}

// Deduct locked stock
func (m *customStocksModel) Deduct(ctx context.Context, session sqlx.Session, skuId, warehouseId uint64, quantity int64) error {
    query := fmt.Sprintf("update %s set `locked` = `locked` - ? where `sku_id` = ? and `warehouse_id` = ? and `locked` >= ?", m.table)
    result, err := session.ExecCtx(ctx, query, quantity, skuId, warehouseId, quantity)
    if err != nil {
        return err
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return ErrNotFound
    }
    return nil
}

// List stocks by warehouse
func (m *customStocksModel) FindManyByWarehouse(ctx context.Context, warehouseId uint64, page, pageSize int32) ([]*Stocks, error) {
    query := fmt.Sprintf("select %s from %s where `warehouse_id` = ? limit ?, ?", stocksRows, m.table)
    var stocks []*Stocks
    err := m.QueryRowsNoCacheCtx(ctx, &stocks, query, warehouseId, (page-1)*pageSize, pageSize)
    return stocks, err
}

// Count stocks by warehouse
func (m *customStocksModel) Count(ctx context.Context, warehouseId uint64) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `warehouse_id` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, warehouseId)
    return count, err
}