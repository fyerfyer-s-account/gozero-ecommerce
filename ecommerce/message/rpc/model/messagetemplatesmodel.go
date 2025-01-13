package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MessageTemplatesModel = (*customMessageTemplatesModel)(nil)

type (
    MessageTemplatesModel interface {
        messageTemplatesModel
        Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
        FindByCode(ctx context.Context, code string) (*MessageTemplates, error)
        ListByTypeAndStatus(ctx context.Context, templateType, status int64, page, pageSize int) ([]*MessageTemplates, error)
        CountByTypeAndStatus(ctx context.Context, templateType, status int64) (int64, error)
        UpdateStatus(ctx context.Context, id uint64, status int64) error
    }

    customMessageTemplatesModel struct {
        *defaultMessageTemplatesModel
    }
)

func NewMessageTemplatesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MessageTemplatesModel {
    return &customMessageTemplatesModel{
        defaultMessageTemplatesModel: newMessageTemplatesModel(conn, c, opts...),
    }
}

func (m *customMessageTemplatesModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
    return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
        return fn(ctx, session)
    })
}

func (m *customMessageTemplatesModel) FindByCode(ctx context.Context, code string) (*MessageTemplates, error) {
    var template MessageTemplates
    templateCodeKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesCodePrefix, code)
    err := m.QueryRowIndexCtx(ctx, &template, templateCodeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (any, error) {
        query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", messageTemplatesRows, m.table)
        if err := conn.QueryRowCtx(ctx, &template, query, code); err != nil {
            return nil, err
        }
        return template.Id, nil
    }, m.queryPrimary)
    switch err {
    case nil:
        return &template, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *customMessageTemplatesModel) ListByTypeAndStatus(ctx context.Context, templateType, status int64, page, pageSize int) ([]*MessageTemplates, error) {
    var templates []*MessageTemplates
    where := "1=1"
    args := make([]interface{}, 0)

    if templateType > 0 {
        where += " AND `type` = ?"
        args = append(args, templateType)
    }
    if status > 0 {
        where += " AND `status` = ?"
        args = append(args, status)
    }

    offset := (page - 1) * pageSize
    query := fmt.Sprintf("select %s from %s where %s order by id desc limit ?, ?", 
        messageTemplatesRows, m.table, where)
    args = append(args, offset, pageSize)

    err := m.QueryRowsNoCacheCtx(ctx, &templates, query, args...)
    return templates, err
}

func (m *customMessageTemplatesModel) CountByTypeAndStatus(ctx context.Context, templateType, status int64) (int64, error) {
    var count int64
    where := "1=1"
    args := make([]interface{}, 0)

    if templateType > 0 {
        where += " AND `type` = ?"
        args = append(args, templateType)
    }
    if status > 0 {
        where += " AND `status` = ?"
        args = append(args, status)
    }

    query := fmt.Sprintf("select count(*) from %s where %s", m.table, where)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
    return count, err
}

func (m *customMessageTemplatesModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
    template, err := m.FindOne(ctx, id)
    if err != nil {
        return err
    }

    templateCodeKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesCodePrefix, template.Code)
    templateIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesIdPrefix, id)
    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, id)
    }, templateCodeKey, templateIdKey)
    return err
}