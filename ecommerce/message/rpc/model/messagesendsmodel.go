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

var _ MessageSendsModel = (*customMessageSendsModel)(nil)

type (
    MessageSendsModel interface {
        messageSendsModel
        Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
        BatchInsert(ctx context.Context, records []*MessageSends) error
        FindByMessageId(ctx context.Context, messageId uint64) ([]*MessageSends, error)
        FindByUserId(ctx context.Context, userId uint64) ([]*MessageSends, error)
        UpdateStatus(ctx context.Context, id uint64, status int64, error string) error
        FindPendingRetry(ctx context.Context, limit int) ([]*MessageSends, error)
        UpdateRetryInfo(ctx context.Context, id uint64, retryCount int64, nextRetryTime time.Time) error
        UpdateSendTime(ctx context.Context, id uint64, sendTime time.Time) error
        FindByStatus(ctx context.Context, status int64, limit int) ([]*MessageSends, error)
    }

    customMessageSendsModel struct {
        *defaultMessageSendsModel
    }
)

// NewMessageSendsModel returns a model for the database table.
func NewMessageSendsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MessageSendsModel {
    return &customMessageSendsModel{
        defaultMessageSendsModel: newMessageSendsModel(conn, c, opts...),
    }
}

// Transaction support
func (m *customMessageSendsModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
    return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
        return fn(ctx, session)
    })
}

// BatchInsert multiple send records
func (m *customMessageSendsModel) BatchInsert(ctx context.Context, records []*MessageSends) error {
    if len(records) == 0 {
        return nil
    }
    values := make([]string, 0, len(records))
    args := make([]interface{}, 0, len(records)*9)
    for _, record := range records {
        values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
        args = append(args, record.MessageId, record.TemplateId, record.UserId,
            record.Channel, record.Status, record.Error, record.RetryCount,
            record.NextRetryTime, record.SendTime)
    }
    query := fmt.Sprintf("insert into %s (%s) values %s",
        m.table, messageSendsRowsExpectAutoSet, strings.Join(values, ","))
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })
    return err
}

// FindByMessageId gets all send records for a message
func (m *customMessageSendsModel) FindByMessageId(ctx context.Context, messageId uint64) ([]*MessageSends, error) {
    var records []*MessageSends
    query := fmt.Sprintf("select %s from %s where `message_id` = ?", messageSendsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, messageId)
    return records, err
}

// FindByUserId gets all send records for a user
func (m *customMessageSendsModel) FindByUserId(ctx context.Context, userId uint64) ([]*MessageSends, error) {
    var records []*MessageSends
    query := fmt.Sprintf("select %s from %s where `user_id` = ?", messageSendsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, userId)
    return records, err
}

// UpdateStatus updates send status and error message
func (m *customMessageSendsModel) UpdateStatus(ctx context.Context, id uint64, status int64, errorStr string) error {
    messageSendsIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageSendsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ?, `error` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, errorStr, id)
    }, messageSendsIdKey)
    return err
}

// FindPendingRetry gets records that need retry
func (m *customMessageSendsModel) FindPendingRetry(ctx context.Context, limit int) ([]*MessageSends, error) {
    var records []*MessageSends
    query := fmt.Sprintf("select %s from %s where `status` = 4 and `retry_count` < 3 and `next_retry_time` <= ? limit ?", 
        messageSendsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, time.Now(), limit)
    return records, err
}

// UpdateRetryInfo updates retry count and next retry time
func (m *customMessageSendsModel) UpdateRetryInfo(ctx context.Context, id uint64, retryCount int64, nextRetryTime time.Time) error {
    messageSendsIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageSendsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `retry_count` = ?, `next_retry_time` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, retryCount, nextRetryTime, id)
    }, messageSendsIdKey)
    return err
}

// UpdateSendTime updates the successful send time
func (m *customMessageSendsModel) UpdateSendTime(ctx context.Context, id uint64, sendTime time.Time) error {
    messageSendsIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageSendsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `send_time` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, sendTime, id)
    }, messageSendsIdKey)
    return err
}

// FindByStatus gets records by status
func (m *customMessageSendsModel) FindByStatus(ctx context.Context, status int64, limit int) ([]*MessageSends, error) {
    var records []*MessageSends
    query := fmt.Sprintf("select %s from %s where `status` = ? limit ?", messageSendsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, status, limit)
    return records, err
}