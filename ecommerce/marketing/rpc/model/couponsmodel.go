package model

import (
    "context"
    "database/sql"
    "fmt"
    "strings"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CouponsModel = (*customCouponsModel)(nil)

type (
    // CouponsModel is an interface to be customized
    CouponsModel interface {
        couponsModel
        Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
        Lock(ctx context.Context, session sqlx.Session, id uint64) error
        Unlock(ctx context.Context, session sqlx.Session, id uint64) error
        IncrReceived(ctx context.Context, id uint64) error
        DecrReceived(ctx context.Context, id uint64) error
        IncrUsed(ctx context.Context, id uint64) error
        FindMany(ctx context.Context, status int32, page, pageSize int32) ([]*Coupons, error)
        Count(ctx context.Context, status int32) (int64, error)
        FindManyByIds(ctx context.Context, ids []uint64) (map[uint64]*Coupons, error)
    }

    customCouponsModel struct {
        *defaultCouponsModel
    }
)

// NewCouponsModel returns a model for the database table.
func NewCouponsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CouponsModel {
    return &customCouponsModel{
        defaultCouponsModel: newCouponsModel(conn, c, opts...),
    }
}

func (m *customCouponsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, fn)
}

// Lock acquires a lock on coupon
func (m *customCouponsModel) Lock(ctx context.Context, session sqlx.Session, id uint64) error {
    query := fmt.Sprintf("select id from %s where `id` = ? for update", m.table)
    var mid uint64
    err := session.QueryRowCtx(ctx, &mid, query, id)
    switch err {
    case nil:
        return nil
    case sqlx.ErrNotFound:
        return ErrNotFound
    default:
        return err
    }
}

// Unlock releases the lock (no-op in MySQL as locks are automatically released)
func (m *customCouponsModel) Unlock(ctx context.Context, session sqlx.Session, id uint64) error {
    return nil
}

// IncrReceived increments received count
func (m *customCouponsModel) IncrReceived(ctx context.Context, id uint64) error {
    key := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `received` = `received` + 1 where `id` = ? and `received` < `total`", m.table)
        return conn.ExecCtx(ctx, query, id)
    }, key)
    return err
}

// DecrReceived decrements received count
func (m *customCouponsModel) DecrReceived(ctx context.Context, id uint64) error {
    key := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `received` = `received` - 1 where `id` = ? and `received` > 0", m.table)
        return conn.ExecCtx(ctx, query, id)
    }, key)
    return err
}

// IncrUsed increments used count
func (m *customCouponsModel) IncrUsed(ctx context.Context, id uint64) error {
    key := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `used` = `used` + 1 where `id` = ? and `used` < `received`", m.table)
        return conn.ExecCtx(ctx, query, id)
    }, key)
    return err
}

// FindMany returns a list of coupons
func (m *customCouponsModel) FindMany(ctx context.Context, status int32, page, pageSize int32) ([]*Coupons, error) {
    var coupons []*Coupons
    query := fmt.Sprintf("select %s from %s where `status` = ? limit ?, ?", couponsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &coupons, query, status, (page-1)*pageSize, pageSize)
    return coupons, err
}

// Count returns total number of coupons
func (m *customCouponsModel) Count(ctx context.Context, status int32) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `status` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, status)
    return count, err
}

// FindManyByIds returns coupons by ids
func (m *customCouponsModel) FindManyByIds(ctx context.Context, ids []uint64) (map[uint64]*Coupons, error) {
    if len(ids) == 0 {
        return make(map[uint64]*Coupons), nil
    }

    // Convert ids to strings for IN clause
    strIds := make([]string, len(ids))
    for i, id := range ids {
        strIds[i] = fmt.Sprint(id)
    }

    var coupons []*Coupons
    query := fmt.Sprintf("select %s from %s where `id` in (%s)", 
        couponsRows, m.table, strings.Join(strIds, ","))
    err := m.QueryRowsNoCacheCtx(ctx, &coupons, query)
    if err != nil {
        return nil, err
    }

    result := make(map[uint64]*Coupons)
    for _, coupon := range coupons {
        result[coupon.Id] = coupon
    }
    return result, nil
}