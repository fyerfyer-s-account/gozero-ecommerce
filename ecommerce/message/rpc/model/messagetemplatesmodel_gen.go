// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	messageTemplatesFieldNames          = builder.RawFieldNames(&MessageTemplates{})
	messageTemplatesRows                = strings.Join(messageTemplatesFieldNames, ",")
	messageTemplatesRowsExpectAutoSet   = strings.Join(stringx.Remove(messageTemplatesFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	messageTemplatesRowsWithPlaceHolder = strings.Join(stringx.Remove(messageTemplatesFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallMessageMessageTemplatesIdPrefix   = "cache:mallMessage:messageTemplates:id:"
	cacheMallMessageMessageTemplatesCodePrefix = "cache:mallMessage:messageTemplates:code:"
)

type (
	messageTemplatesModel interface {
		Insert(ctx context.Context, data *MessageTemplates) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*MessageTemplates, error)
		FindOneByCode(ctx context.Context, code string) (*MessageTemplates, error)
		Update(ctx context.Context, data *MessageTemplates) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultMessageTemplatesModel struct {
		sqlc.CachedConn
		table string
	}

	MessageTemplates struct {
		Id              uint64         `db:"id"`               // æ¨¡æ¿ID
		Code            string         `db:"code"`             // æ¨¡æ¿ç¼–ç 
		Name            string         `db:"name"`             // æ¨¡æ¿åç§°
		TitleTemplate   string         `db:"title_template"`   // æ ‡é¢˜æ¨¡æ¿
		ContentTemplate string         `db:"content_template"` // å†…å®¹æ¨¡æ¿
		Type            int64          `db:"type"`             // æ¶ˆæ¯ç±»åž‹
		Channels        string         `db:"channels"`         // å‘é€æ¸ é“
		Config          sql.NullString `db:"config"`           // æ¸ é“é…ç½®
		Status          int64          `db:"status"`           // çŠ¶æ€ 1:å¯ç”¨ 2:ç¦ç”¨
		CreatedAt       time.Time      `db:"created_at"`       // åˆ›å»ºæ—¶é—´
		UpdatedAt       time.Time      `db:"updated_at"`       // æ›´æ–°æ—¶é—´
	}
)

func newMessageTemplatesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultMessageTemplatesModel {
	return &defaultMessageTemplatesModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`message_templates`",
	}
}

func (m *defaultMessageTemplatesModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallMessageMessageTemplatesCodeKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesCodePrefix, data.Code)
	mallMessageMessageTemplatesIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mallMessageMessageTemplatesCodeKey, mallMessageMessageTemplatesIdKey)
	return err
}

func (m *defaultMessageTemplatesModel) FindOne(ctx context.Context, id uint64) (*MessageTemplates, error) {
	mallMessageMessageTemplatesIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesIdPrefix, id)
	var resp MessageTemplates
	err := m.QueryRowCtx(ctx, &resp, mallMessageMessageTemplatesIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messageTemplatesRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMessageTemplatesModel) FindOneByCode(ctx context.Context, code string) (*MessageTemplates, error) {
	mallMessageMessageTemplatesCodeKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesCodePrefix, code)
	var resp MessageTemplates
	err := m.QueryRowIndexCtx(ctx, &resp, mallMessageMessageTemplatesCodeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", messageTemplatesRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, code); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMessageTemplatesModel) Insert(ctx context.Context, data *MessageTemplates) (sql.Result, error) {
	mallMessageMessageTemplatesCodeKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesCodePrefix, data.Code)
	mallMessageMessageTemplatesIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, messageTemplatesRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Code, data.Name, data.TitleTemplate, data.ContentTemplate, data.Type, data.Channels, data.Config, data.Status)
	}, mallMessageMessageTemplatesCodeKey, mallMessageMessageTemplatesIdKey)
	return ret, err
}

func (m *defaultMessageTemplatesModel) Update(ctx context.Context, newData *MessageTemplates) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	mallMessageMessageTemplatesCodeKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesCodePrefix, data.Code)
	mallMessageMessageTemplatesIdKey := fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, messageTemplatesRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.Code, newData.Name, newData.TitleTemplate, newData.ContentTemplate, newData.Type, newData.Channels, newData.Config, newData.Status, newData.Id)
	}, mallMessageMessageTemplatesCodeKey, mallMessageMessageTemplatesIdKey)
	return err
}

func (m *defaultMessageTemplatesModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallMessageMessageTemplatesIdPrefix, primary)
}

func (m *defaultMessageTemplatesModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messageTemplatesRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultMessageTemplatesModel) tableName() string {
	return m.table
}
