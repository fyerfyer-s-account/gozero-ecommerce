package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		WithSession(session sqlx.Session) UsersModel
		FindOneByPhoneOrEmail(ctx context.Context, account string) (*Users, error)
		UpdatePassword(ctx context.Context, id uint64, password string) error
		UpdateStatus(ctx context.Context, id uint64, status int32) error
		UpdateProfile(ctx context.Context, id uint64, nickname string, gender string) error
		UpdateOnline(ctx context.Context, id uint64, online int32) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) WithSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUsersModel) FindOneByPhoneOrEmail(ctx context.Context, account string) (*Users, error) {
	query := fmt.Sprintf("select %s from %s where `username` = ? or `email` = ? limit 1", usersRows, m.table)
	var resp Users
	err := m.conn.QueryRowCtx(ctx, &resp, query, account, account)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) UpdatePassword(ctx context.Context, id uint64, password string) error {
	query := `update ` + m.table + ` set password = ? where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, password, id)
	return err
}

func (m *customUsersModel) UpdateStatus(ctx context.Context, id uint64, status int32) error {
	query := `update ` + m.table + ` set status = ? where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, status, id)
	return err
}

func (m *customUsersModel) UpdateProfile(ctx context.Context, id uint64, nickname string, gender string) error {
	query := `update ` + m.table + ` set nickname = ?, gender = ? where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, nickname, gender, id)
	return err
}

func (m *customUsersModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customUsersModel) UpdateOnline(ctx context.Context, id uint64, online int32) error {
	query := `update ` + m.table + ` set online = ? where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, online, id)
	return err
}
