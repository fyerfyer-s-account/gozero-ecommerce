package model

import (
	"context"
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
		UpdateProfile(ctx context.Context, id uint64, nickname string, avatar string, gender int32) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) FindOneByPhoneOrEmail(ctx context.Context, account string) (*Users, error) {
	var user Users
	query := `select ` + usersRows + ` from ` + m.table + ` where phone = ? or email = ? limit 1`
	err := m.conn.QueryRowCtx(ctx, &user, query, account, account)
	if err != nil {
		return nil, err
	}
	return &user, nil
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

func (m *customUsersModel) UpdateProfile(ctx context.Context, id uint64, nickname string, avatar string, gender int32) error {
	query := `update ` + m.table + ` set nickname = ?, avatar = ?, gender = ? where id = ?`
	_, err := m.conn.ExecCtx(ctx, query, nickname, avatar, gender, id)
	return err
}

func (m *customUsersModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customUsersModel) WithSession(session sqlx.Session) UsersModel {
	return &customUsersModel{
		defaultUsersModel: &defaultUsersModel{
			conn:  sqlx.NewSqlConnFromSession(session),
			table: m.table,
		},
	}
}
