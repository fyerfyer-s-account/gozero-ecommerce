package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StockLocksModel = (*customStockLocksModel)(nil)

type (
	// StockLocksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockLocksModel.
	StockLocksModel interface {
		stockLocksModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
        FindByOrderNo(ctx context.Context, orderNo string) ([]*StockLocks, error)
        BatchInsert(ctx context.Context, data []*StockLocks) error
        UpdateStatus(ctx context.Context, orderNo string, oldStatus, newStatus int64) error
        FindAndLockByOrderNo(ctx context.Context, session sqlx.Session, orderNo string) ([]*StockLocks, error)
        DeleteByOrderNo(ctx context.Context, orderNo string) error
	}

	customStockLocksModel struct {
		*defaultStockLocksModel
	}
)

// NewStockLocksModel returns a model for the database table.
func NewStockLocksModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StockLocksModel {
	return &customStockLocksModel{
		defaultStockLocksModel: newStockLocksModel(conn, c, opts...),
	}
}

// Transaction support
func (m *customStockLocksModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
        return fn(ctx, session)
    })
}

// Find by order number
func (m *customStockLocksModel) FindByOrderNo(ctx context.Context, orderNo string) ([]*StockLocks, error) {
    var locks []*StockLocks
    query := fmt.Sprintf("select %s from %s where `order_no` = ?", stockLocksRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &locks, query, orderNo)
    return locks, err
}

// Batch insert locks
func (m *customStockLocksModel) BatchInsert(ctx context.Context, data []*StockLocks) error {
    if len(data) == 0 {
        return nil
    }
    
    values := make([]string, 0, len(data))
    args := make([]interface{}, 0, len(data)*5)
    for _, lock := range data {
        values = append(values, "(?, ?, ?, ?, ?)")
        args = append(args, lock.OrderNo, lock.SkuId, lock.WarehouseId, lock.Quantity, lock.Status)
    }
    
    query := fmt.Sprintf("insert into %s (%s) values %s",
        m.table, stockLocksRowsExpectAutoSet, strings.Join(values, ","))
    
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })
    return err
}

// Update lock status
func (m *customStockLocksModel) UpdateStatus(ctx context.Context, orderNo string, oldStatus, newStatus int64) error {
    query := fmt.Sprintf("update %s set `status` = ? where `order_no` = ? and `status` = ?", m.table)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, newStatus, orderNo, oldStatus)
    })
    return err
}

// Find and lock records for update
func (m *customStockLocksModel) FindAndLockByOrderNo(ctx context.Context, session sqlx.Session, orderNo string) ([]*StockLocks, error) {
    var locks []*StockLocks
    query := fmt.Sprintf("select %s from %s where `order_no` = ? for update", stockLocksRows, m.table)
    err := session.QueryRowsCtx(ctx, &locks, query, orderNo)
    return locks, err
}

// Delete by order number
func (m *customStockLocksModel) DeleteByOrderNo(ctx context.Context, orderNo string) error {
    // First get all locks to invalidate cache
    locks, err := m.FindByOrderNo(ctx, orderNo)
    if err != nil {
        return err
    }
    
    keys := make([]string, 0, len(locks))
    for _, lock := range locks {
        keys = append(keys, fmt.Sprintf("%s%v", cacheMallInventoryStockLocksIdPrefix, lock.Id))
    }
    
    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("delete from %s where `order_no` = ?", m.table)
        return conn.ExecCtx(ctx, query, orderNo)
    }, keys...)
    
    return err
}