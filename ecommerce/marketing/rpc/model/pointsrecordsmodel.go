package model

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PointsRecordsModel = (*customPointsRecordsModel)(nil)

type (
    PointsRecordsModel interface {
        pointsRecordsModel
        Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
        FindByUserId(ctx context.Context, userId int64, page, pageSize int32) ([]*PointsRecords, error)
        BatchInsert(ctx context.Context, records []*PointsRecords) error
        CountByUserId(ctx context.Context, userId int64) (int64, error)
        SumPointsByUserId(ctx context.Context, userId int64) (int64, error)
        FindByDateRange(ctx context.Context, userId int64, startTime, endTime time.Time) ([]*PointsRecords, error)
    }

    customPointsRecordsModel struct {
        *defaultPointsRecordsModel
    }
)

func NewPointsRecordsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PointsRecordsModel {
    return &customPointsRecordsModel{
        defaultPointsRecordsModel: newPointsRecordsModel(conn, c, opts...),
    }
}

func (m *customPointsRecordsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, fn)
}

func (m *customPointsRecordsModel) FindByUserId(ctx context.Context, userId int64, page, pageSize int32) ([]*PointsRecords, error) {
    var records []*PointsRecords
    query := fmt.Sprintf("select %s from %s where `user_id` = ? order by id desc limit ?, ?", pointsRecordsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, userId, (page-1)*pageSize, pageSize)
    return records, err
}

func (m *customPointsRecordsModel) BatchInsert(ctx context.Context, records []*PointsRecords) error {
    if len(records) == 0 {
        return nil
    }

    values := make([]string, 0, len(records))
    args := make([]interface{}, 0, len(records)*6)
    
    for _, record := range records {
        values = append(values, "(?, ?, ?, ?, ?, ?)")
        args = append(args, record.UserId, record.Points, record.Type, 
            record.Source, record.Remark, record.CreatedAt)
    }

    query := fmt.Sprintf("insert into %s (user_id, points, type, source, remark, created_at) values %s",
        m.table, strings.Join(values, ","))
    
    _, err := m.ExecNoCacheCtx(ctx, query, args...)
    return err
}

func (m *customPointsRecordsModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, userId)
    return count, err
}

func (m *customPointsRecordsModel) SumPointsByUserId(ctx context.Context, userId int64) (int64, error) {
    var sum int64
    query := fmt.Sprintf("select COALESCE(sum(case when type = 1 then points else -points end), 0) from %s where `user_id` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &sum, query, userId)
    return sum, err
}

func (m *customPointsRecordsModel) FindByDateRange(ctx context.Context, userId int64, startTime, endTime time.Time) ([]*PointsRecords, error) {
    var records []*PointsRecords
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and `created_at` between ? and ? order by id desc",
        pointsRecordsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, userId, startTime, endTime)
    return records, err
}