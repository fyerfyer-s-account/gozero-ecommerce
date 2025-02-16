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
	adminsFieldNames          = builder.RawFieldNames(&Admins{})
	adminsRows                = strings.Join(adminsFieldNames, ",")
	adminsRowsExpectAutoSet   = strings.Join(stringx.Remove(adminsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	adminsRowsWithPlaceHolder = strings.Join(stringx.Remove(adminsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	adminsModel interface {
		Insert(ctx context.Context, data *Admins) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Admins, error)
		FindOneByUserId(ctx context.Context, userId uint64) (*Admins, error)
		Update(ctx context.Context, data *Admins) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultAdminsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Admins struct {
		Id        uint64    `db:"id"`         // ç®¡ç†å‘˜ID
		UserId    uint64    `db:"user_id"`    // ç”¨æˆ·ID
		CreatedAt time.Time `db:"created_at"` // åˆ›å»ºæ—¶é—´
	}
)

func newAdminsModel(conn sqlx.SqlConn) *defaultAdminsModel {
	return &defaultAdminsModel{
		conn:  conn,
		table: "`admins`",
	}
}

func (m *defaultAdminsModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultAdminsModel) FindOne(ctx context.Context, id uint64) (*Admins, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", adminsRows, m.table)
	var resp Admins
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

func (m *defaultAdminsModel) FindOneByUserId(ctx context.Context, userId uint64) (*Admins, error) {
	var resp Admins
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", adminsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAdminsModel) Insert(ctx context.Context, data *Admins) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?)", m.table, adminsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId)
	return ret, err
}

func (m *defaultAdminsModel) Update(ctx context.Context, newData *Admins) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, adminsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.UserId, newData.Id)
	return err
}

func (m *defaultAdminsModel) tableName() string {
	return m.table
}
