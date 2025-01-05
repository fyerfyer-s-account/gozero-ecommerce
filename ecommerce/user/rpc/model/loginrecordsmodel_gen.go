// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.4

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
	loginRecordsFieldNames          = builder.RawFieldNames(&LoginRecords{})
	loginRecordsRows                = strings.Join(loginRecordsFieldNames, ",")
	loginRecordsRowsExpectAutoSet   = strings.Join(stringx.Remove(loginRecordsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	loginRecordsRowsWithPlaceHolder = strings.Join(stringx.Remove(loginRecordsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	loginRecordsModel interface {
		Insert(ctx context.Context, data *LoginRecords) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*LoginRecords, error)
		Update(ctx context.Context, data *LoginRecords) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultLoginRecordsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	LoginRecords struct {
		Id            uint64         `db:"id"`             // è®°å½•ID
		UserId        uint64         `db:"user_id"`        // ç”¨æˆ·ID
		LoginIp       string         `db:"login_ip"`       // ç™»å½•IP
		LoginLocation sql.NullString `db:"login_location"` // ç™»å½•åœ°ç‚¹
		DeviceType    sql.NullString `db:"device_type"`    // è®¾å¤‡ç±»åž‹
		CreatedAt     time.Time      `db:"created_at"`     // åˆ›å»ºæ—¶é—´
	}
)

func newLoginRecordsModel(conn sqlx.SqlConn) *defaultLoginRecordsModel {
	return &defaultLoginRecordsModel{
		conn:  conn,
		table: "`login_records`",
	}
}

func (m *defaultLoginRecordsModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultLoginRecordsModel) FindOne(ctx context.Context, id uint64) (*LoginRecords, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", loginRecordsRows, m.table)
	var resp LoginRecords
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

func (m *defaultLoginRecordsModel) Insert(ctx context.Context, data *LoginRecords) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, loginRecordsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.LoginIp, data.LoginLocation, data.DeviceType)
	return ret, err
}

func (m *defaultLoginRecordsModel) Update(ctx context.Context, data *LoginRecords) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, loginRecordsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.LoginIp, data.LoginLocation, data.DeviceType, data.Id)
	return err
}

func (m *defaultLoginRecordsModel) tableName() string {
	return m.table
}
