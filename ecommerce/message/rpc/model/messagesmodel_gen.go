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
	messagesFieldNames          = builder.RawFieldNames(&Messages{})
	messagesRows                = strings.Join(messagesFieldNames, ",")
	messagesRowsExpectAutoSet   = strings.Join(stringx.Remove(messagesFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	messagesRowsWithPlaceHolder = strings.Join(stringx.Remove(messagesFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	messagesModel interface {
		Insert(ctx context.Context, data *Messages) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Messages, error)
		Update(ctx context.Context, data *Messages) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultMessagesModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Messages struct {
		Id          uint64         `db:"id"`           // æ¶ˆæ¯ID
		UserId      uint64         `db:"user_id"`      // ç”¨æˆ·ID
		Title       string         `db:"title"`        // æ¶ˆæ¯æ ‡é¢˜
		Content     string         `db:"content"`      // æ¶ˆæ¯å†…å®¹
		Type        int64          `db:"type"`         // æ¶ˆæ¯ç±»åž‹ 1:ç³»ç»Ÿé€šçŸ¥ 2:è®¢å•æ¶ˆæ¯ 3:æ´»åŠ¨æ¶ˆæ¯ 4:ç‰©æµæ¶ˆæ¯
		SendChannel int64          `db:"send_channel"` // å‘é€æ¸ é“ 1:ç«™å†…ä¿¡ 2:çŸ­ä¿¡ 3:é‚®ä»¶ 4:APPæŽ¨é€
		ExtraData   sql.NullString `db:"extra_data"`   // é¢å¤–æ•°æ®
		IsRead      int64          `db:"is_read"`      // æ˜¯å¦å·²è¯»
		ReadTime    sql.NullTime   `db:"read_time"`    // é˜…è¯»æ—¶é—´
		CreatedAt   time.Time      `db:"created_at"`   // åˆ›å»ºæ—¶é—´
	}
)

func newMessagesModel(conn sqlx.SqlConn) *defaultMessagesModel {
	return &defaultMessagesModel{
		conn:  conn,
		table: "`messages`",
	}
}

func (m *defaultMessagesModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultMessagesModel) FindOne(ctx context.Context, id uint64) (*Messages, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", messagesRows, m.table)
	var resp Messages
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

func (m *defaultMessagesModel) Insert(ctx context.Context, data *Messages) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, messagesRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Title, data.Content, data.Type, data.SendChannel, data.ExtraData, data.IsRead, data.ReadTime)
	return ret, err
}

func (m *defaultMessagesModel) Update(ctx context.Context, data *Messages) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, messagesRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Title, data.Content, data.Type, data.SendChannel, data.ExtraData, data.IsRead, data.ReadTime, data.Id)
	return err
}

func (m *defaultMessagesModel) tableName() string {
	return m.table
}
