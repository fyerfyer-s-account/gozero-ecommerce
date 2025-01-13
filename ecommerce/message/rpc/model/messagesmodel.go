package model

import (
    "context"
    "database/sql"
    "fmt"
    "strings"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessagesModel = (*customMessagesModel)(nil)

type (
    MessagesModel interface {
        messagesModel
        Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
        BatchInsert(ctx context.Context, messages []*Messages) error
        FindByUserId(ctx context.Context, userId uint64, messageType int64, unreadOnly bool, page, pageSize int) ([]*Messages, error)
        CountByUserId(ctx context.Context, userId uint64, messageType int64, unreadOnly bool) (int64, error)
        UpdateReadStatus(ctx context.Context, id, userId uint64) error
        DeleteByUserMessage(ctx context.Context, id, userId uint64) error
    }

    customMessagesModel struct {
        *defaultMessagesModel
    }
)

func NewMessagesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MessagesModel {
    return &customMessagesModel{
        defaultMessagesModel: newMessagesModel(conn, c, opts...),
    }
}

func (m *customMessagesModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
    return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
        return fn(ctx, session)
    })
}

func (m *customMessagesModel) BatchInsert(ctx context.Context, messages []*Messages) error {
    if len(messages) == 0 {
        return nil
    }
    
    values := make([]string, 0, len(messages))
    args := make([]interface{}, 0, len(messages)*8)
    for _, msg := range messages {
        values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?)")
        args = append(args, msg.UserId, msg.Title, msg.Content, msg.Type,
            msg.SendChannel, msg.ExtraData, msg.IsRead, msg.ReadTime)
    }
    
    query := fmt.Sprintf("insert into %s (%s) values %s",
        m.table, messagesRowsExpectAutoSet, strings.Join(values, ","))
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })
    return err
}

func (m *customMessagesModel) FindByUserId(ctx context.Context, userId uint64, messageType int64, unreadOnly bool, page, pageSize int) ([]*Messages, error) {
    var messages []*Messages
    where := "`user_id` = ?"
    args := []interface{}{userId}
    
    if messageType > 0 {
        where += " AND `type` = ?"
        args = append(args, messageType)
    }
    if unreadOnly {
        where += " AND `is_read` = 0"
    }
    
    offset := (page - 1) * pageSize
    query := fmt.Sprintf("select %s from %s where %s order by created_at desc limit ?, ?",
        messagesRows, m.table, where)
    args = append(args, offset, pageSize)
    
    err := m.QueryRowsNoCacheCtx(ctx, &messages, query, args...)
    return messages, err
}

func (m *customMessagesModel) CountByUserId(ctx context.Context, userId uint64, messageType int64, unreadOnly bool) (int64, error) {
    var count int64
    where := "`user_id` = ?"
    args := []interface{}{userId}
    
    if messageType > 0 {
        where += " AND `type` = ?"
        args = append(args, messageType)
    }
    if unreadOnly {
        where += " AND `is_read` = 0"
    }
    
    query := fmt.Sprintf("select count(*) from %s where %s", m.table, where)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
    return count, err
}

func (m *customMessagesModel) UpdateReadStatus(ctx context.Context, id, userId uint64) error {
    messageIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessagesIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `is_read` = 1, `read_time` = now() where `id` = ? and `user_id` = ? and `is_read` = 0", m.table)
        return conn.ExecCtx(ctx, query, id, userId)
    }, messageIdKey)
    return err
}

func (m *customMessagesModel) DeleteByUserMessage(ctx context.Context, id, userId uint64) error {
    messageIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessagesIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("delete from %s where `id` = ? and `user_id` = ?", m.table)
        return conn.ExecCtx(ctx, query, id, userId)
    }, messageIdKey)
    return err
}