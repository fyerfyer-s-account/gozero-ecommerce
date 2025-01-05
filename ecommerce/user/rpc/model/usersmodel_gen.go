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
	usersFieldNames          = builder.RawFieldNames(&Users{})
	usersRows                = strings.Join(usersFieldNames, ",")
	usersRowsExpectAutoSet   = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	usersRowsWithPlaceHolder = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	usersModel interface {
		Insert(ctx context.Context, data *Users) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Users, error)
		FindOneByEmail(ctx context.Context, email sql.NullString) (*Users, error)
		FindOneByPhone(ctx context.Context, phone sql.NullString) (*Users, error)
		FindOneByUsername(ctx context.Context, username string) (*Users, error)
		Update(ctx context.Context, data *Users) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultUsersModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Users struct {
		Id          uint64         `db:"id"`           // ç”¨æˆ·ID
		Username    string         `db:"username"`     // ç”¨æˆ·å
		Password    string         `db:"password"`     // å¯†ç 
		Phone       sql.NullString `db:"phone"`        // æ‰‹æœºå·
		Email       sql.NullString `db:"email"`        // é‚®ç®±
		Nickname    sql.NullString `db:"nickname"`     // æ˜µç§°
		Avatar      sql.NullString `db:"avatar"`       // å¤´åƒURL
		Gender      string         `db:"gender"`       // æ€§åˆ«
		MemberLevel int64          `db:"member_level"` // ä¼šå‘˜ç­‰çº§
		Status      int64          `db:"status"`       // çŠ¶æ€ 0:ç¦ç”¨ 1:å¯ç”¨
		Online      int64          `db:"online"`       // çŠ¶æ€ 0:ç¦»çº¿ 1:åœ¨çº¿
		CreatedAt   time.Time      `db:"created_at"`   // åˆ›å»ºæ—¶é—´
		UpdatedAt   time.Time      `db:"updated_at"`   // æ›´æ–°æ—¶é—´
	}
)

func newUsersModel(conn sqlx.SqlConn) *defaultUsersModel {
	return &defaultUsersModel{
		conn:  conn,
		table: "`users`",
	}
}

func (m *defaultUsersModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultUsersModel) FindOne(ctx context.Context, id uint64) (*Users, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
	var resp Users
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

func (m *defaultUsersModel) FindOneByEmail(ctx context.Context, email sql.NullString) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUsersModel) FindOneByPhone(ctx context.Context, phone sql.NullString) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUsersModel) FindOneByUsername(ctx context.Context, username string) (*Users, error) {
	var resp Users
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, usersRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Username, data.Password, data.Phone, data.Email, data.Nickname, data.Avatar, data.Gender, data.MemberLevel, data.Status, data.Online)
	return ret, err
}

func (m *defaultUsersModel) Update(ctx context.Context, newData *Users) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, usersRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Username, newData.Password, newData.Phone, newData.Email, newData.Nickname, newData.Avatar, newData.Gender, newData.MemberLevel, newData.Status, newData.Online, newData.Id)
	return err
}

func (m *defaultUsersModel) tableName() string {
	return m.table
}
