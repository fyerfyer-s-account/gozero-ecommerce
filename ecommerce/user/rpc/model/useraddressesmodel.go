package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserAddressesModel = (*customUserAddressesModel)(nil)

type (
	UserAddressesModel interface {
		userAddressesModel
		WithSession(session sqlx.Session) UserAddressesModel
		FindByUserId(ctx context.Context, userId uint64) ([]*UserAddresses, error)
		SetDefault(ctx context.Context, userId, addressId uint64) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
	}

	customUserAddressesModel struct {
		*defaultUserAddressesModel
	}
)

func NewUserAddressesModel(conn sqlx.SqlConn) UserAddressesModel {
	return &customUserAddressesModel{
		defaultUserAddressesModel: newUserAddressesModel(conn),
	}
}

func (m *customUserAddressesModel) FindByUserId(ctx context.Context, userId uint64) ([]*UserAddresses, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", userAddressesRows, m.table)
	var resp []*UserAddresses
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customUserAddressesModel) SetDefault(ctx context.Context, userId, addressId uint64) error {
	return m.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		// Clear previous default address
		_, err := session.ExecCtx(ctx, fmt.Sprintf("update %s set `is_default` = 0 where `user_id` = ?", m.table), userId)
		if err != nil {
			return err
		}

		// Set new default address
		_, err = session.ExecCtx(ctx, fmt.Sprintf("update %s set `is_default` = 1 where `id` = ? and `user_id` = ?", m.table), addressId, userId)
		return err
	})
}

func (m *customUserAddressesModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customUserAddressesModel) WithSession(session sqlx.Session) UserAddressesModel {
	return &customUserAddressesModel{
		defaultUserAddressesModel: &defaultUserAddressesModel{
			conn:  sqlx.NewSqlConnFromSession(session),
			table: m.table,
		},
	}
}
