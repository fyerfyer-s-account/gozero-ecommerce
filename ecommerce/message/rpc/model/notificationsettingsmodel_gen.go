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
	notificationSettingsFieldNames          = builder.RawFieldNames(&NotificationSettings{})
	notificationSettingsRows                = strings.Join(notificationSettingsFieldNames, ",")
	notificationSettingsRowsExpectAutoSet   = strings.Join(stringx.Remove(notificationSettingsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	notificationSettingsRowsWithPlaceHolder = strings.Join(stringx.Remove(notificationSettingsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	notificationSettingsModel interface {
		Insert(ctx context.Context, data *NotificationSettings) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*NotificationSettings, error)
		FindOneByUserIdTypeChannel(ctx context.Context, userId uint64, tp int64, channel int64) (*NotificationSettings, error)
		Update(ctx context.Context, data *NotificationSettings) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultNotificationSettingsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	NotificationSettings struct {
		Id        uint64    `db:"id"`         // è®¾ç½®ID
		UserId    uint64    `db:"user_id"`    // ç”¨æˆ·ID
		Type      int64     `db:"type"`       // æ¶ˆæ¯ç±»åž‹
		Channel   int64     `db:"channel"`    // é€šçŸ¥æ¸ é“
		IsEnabled int64     `db:"is_enabled"` // æ˜¯å¦å¯ç”¨
		CreatedAt time.Time `db:"created_at"` // åˆ›å»ºæ—¶é—´
		UpdatedAt time.Time `db:"updated_at"` // æ›´æ–°æ—¶é—´
	}
)

func newNotificationSettingsModel(conn sqlx.SqlConn) *defaultNotificationSettingsModel {
	return &defaultNotificationSettingsModel{
		conn:  conn,
		table: "`notification_settings`",
	}
}

func (m *defaultNotificationSettingsModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultNotificationSettingsModel) FindOne(ctx context.Context, id uint64) (*NotificationSettings, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", notificationSettingsRows, m.table)
	var resp NotificationSettings
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

func (m *defaultNotificationSettingsModel) FindOneByUserIdTypeChannel(ctx context.Context, userId uint64, tp int64, channel int64) (*NotificationSettings, error) {
	var resp NotificationSettings
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `type` = ? and `channel` = ? limit 1", notificationSettingsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, tp, channel)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNotificationSettingsModel) Insert(ctx context.Context, data *NotificationSettings) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, notificationSettingsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Type, data.Channel, data.IsEnabled)
	return ret, err
}

func (m *defaultNotificationSettingsModel) Update(ctx context.Context, newData *NotificationSettings) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, notificationSettingsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.UserId, newData.Type, newData.Channel, newData.IsEnabled, newData.Id)
	return err
}

func (m *defaultNotificationSettingsModel) tableName() string {
	return m.table
}