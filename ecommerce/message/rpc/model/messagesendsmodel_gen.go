// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	messageSendsFieldNames          = builder.RawFieldNames(&MessageSends{})
	messageSendsRows                = strings.Join(messageSendsFieldNames, ",")
	messageSendsRowsExpectAutoSet   = strings.Join(stringx.Remove(messageSendsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	messageSendsRowsWithPlaceHolder = strings.Join(stringx.Remove(messageSendsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	messageSendsModel interface {
		Insert(ctx context.Context, data *MessageSends) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*MessageSends, error)
		Update(ctx context.Context, data *MessageSends) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultMessageSendsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	MessageSends struct {
		Id            uint64         `db:"id"`              // è®°å½•ID
		MessageId     uint64         `db:"message_id"`      // æ¶ˆæ¯ID
		TemplateId    sql.NullInt64  `db:"template_id"`     // æ¨¡æ¿ID
		UserId        uint64         `db:"user_id"`         // ç”¨æˆ·ID
		Channel       int64          `db:"channel"`         // å‘é€æ¸ é“
		Status        int64          `db:"status"`          // å‘é€çŠ¶æ€ 1:å¾…å‘é€ 2:å‘é€ä¸­ 3:å‘é€æˆåŠŸ 4:å‘é€å¤±è´¥
		Error         sql.NullString `db:"error"`           // é”™è¯¯ä¿¡æ¯
		RetryCount    int64          `db:"retry_count"`     // é‡è¯•æ¬¡æ•°
		NextRetryTime sql.NullTime   `db:"next_retry_time"` // ä¸‹æ¬¡é‡è¯•æ—¶é—´
		SendTime      sql.NullTime   `db:"send_time"`       // å‘é€æ—¶é—´
		CreatedAt     time.Time      `db:"created_at"`      // åˆ›å»ºæ—¶é—´
		UpdatedAt     time.Time      `db:"updated_at"`      // æ›´æ–°æ—¶é—´
	}
)

func newMessageSendsModel(conn sqlx.SqlConn) *defaultMessageSendsModel {
	return &defaultMessageSendsModel{
		conn:  conn,
		table: "`message_sends`",
	}
}

func (m *defaultMessageSendsModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMessageSendsModel) FindOne(ctx context.Context, id uint64) (*MessageSends, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messageSendsRows, m.table)
	var resp MessageSends
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMessageSendsModel) Insert(ctx context.Context, data *MessageSends) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, messageSendsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.MessageId, data.TemplateId, data.UserId, data.Channel, data.Status, data.Error, data.RetryCount, data.NextRetryTime, data.SendTime)
	return ret, err
}

func (m *defaultMessageSendsModel) Update(ctx context.Context, data *MessageSends) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, messageSendsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.MessageId, data.TemplateId, data.UserId, data.Channel, data.Status, data.Error, data.RetryCount, data.NextRetryTime, data.SendTime, data.Id)
	return err
}

func (m *defaultMessageSendsModel) tableName() string {
	return m.table
}
