package model

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
    PromotionsModel interface {
        promotionsModel
        Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
        FindByStatus(ctx context.Context, status int32, page, pageSize int32) ([]*Promotions, error)
        FindActive(ctx context.Context) ([]*Promotions, error)
        UpdateStatus(ctx context.Context, id uint64, status int32) error
        FindByDateRange(ctx context.Context, startTime, endTime time.Time) ([]*Promotions, error)
        Count(ctx context.Context, status int32) (int64, error)
        BatchInsert(ctx context.Context, data []*Promotions) error
    }

    customPromotionsModel struct {
        *defaultPromotionsModel
    }
)

func NewPromotionsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PromotionsModel {
    return &customPromotionsModel{
        defaultPromotionsModel: newPromotionsModel(conn, c, opts...),
    }
}

func (m *customPromotionsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, fn)
}

func (m *customPromotionsModel) FindByStatus(ctx context.Context, status int32, page, pageSize int32) ([]*Promotions, error) {
    var resp []*Promotions
    query := fmt.Sprintf("select %s from %s where `status` = ? limit ?, ?", promotionsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status, (page-1)*pageSize, pageSize)
    return resp, err
}

func (m *customPromotionsModel) FindActive(ctx context.Context) ([]*Promotions, error) {
    var resp []*Promotions
    now := time.Now()
    query := fmt.Sprintf("select %s from %s where `status` = 1 and `start_time` <= ? and `end_time` >= ?", 
        promotionsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, now, now)
    return resp, err
}

func (m *customPromotionsModel) UpdateStatus(ctx context.Context, id uint64, status int32) error {
    mallMarketingPromotionsIdKey := fmt.Sprintf("%s%v", cacheMallMarketingPromotionsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, id)
    }, mallMarketingPromotionsIdKey)
    return err
}

func (m *customPromotionsModel) FindByDateRange(ctx context.Context, startTime, endTime time.Time) ([]*Promotions, error) {
    var resp []*Promotions
    query := fmt.Sprintf("select %s from %s where (`start_time` between ? and ?) or (`end_time` between ? and ?)", 
        promotionsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, startTime, endTime, startTime, endTime)
    return resp, err
}

func (m *customPromotionsModel) Count(ctx context.Context, status int32) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `status` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, status)
    return count, err
}

func (m *customPromotionsModel) BatchInsert(ctx context.Context, data []*Promotions) error {
    if len(data) == 0 {
        return nil
    }

    values := make([]string, 0, len(data))
    args := make([]interface{}, 0, len(data)*6)
    
    for _, promo := range data {
        values = append(values, "(?, ?, ?, ?, ?, ?)")
        args = append(args, promo.Name, promo.Type, promo.Rules, 
            promo.Status, promo.StartTime, promo.EndTime)
    }

    query := fmt.Sprintf("insert into %s (%s) values %s", 
        m.table, promotionsRowsExpectAutoSet, strings.Join(values, ","))
    
    _, err := m.ExecNoCacheCtx(ctx, query, args...)
    return err
}