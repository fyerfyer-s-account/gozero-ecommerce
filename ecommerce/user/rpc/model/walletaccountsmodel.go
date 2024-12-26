package model

import (
	"context"
	"database/sql"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WalletAccountsModel = (*customWalletAccountsModel)(nil)

type (
	WalletAccountsModel interface {
		walletAccountsModel
		withSession(session sqlx.Session) WalletAccountsModel
		UpdateBalance(ctx context.Context, userId uint64, amount float64) error
		FreezeAmount(ctx context.Context, userId uint64, amount float64) error
		UnfreezeAmount(ctx context.Context, userId uint64, amount float64) error
		UpdatePayPassword(ctx context.Context, userId uint64, password string) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
	}

	customWalletAccountsModel struct {
		*defaultWalletAccountsModel
	}
)

func NewWalletAccountsModel(conn sqlx.SqlConn) WalletAccountsModel {
	return &customWalletAccountsModel{
		defaultWalletAccountsModel: newWalletAccountsModel(conn),
	}
}

func (m *customWalletAccountsModel) UpdateBalance(ctx context.Context, userId uint64, amount float64) error {
	query := "update wallet_accounts set balance = balance + ? where user_id = ? and balance + ? >= 0"
	result, err := m.conn.ExecCtx(ctx, query, amount, userId, amount)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (m *customWalletAccountsModel) FreezeAmount(ctx context.Context, userId uint64, amount float64) error {
	query := "update wallet_accounts set balance = balance - ?, frozen_amount = frozen_amount + ? where user_id = ? and balance >= ?"
	result, err := m.conn.ExecCtx(ctx, query, amount, amount, userId, amount)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return zeroerr.ErrInsufficientBalance
	}
	return nil
}

func (m *customWalletAccountsModel) UnfreezeAmount(ctx context.Context, userId uint64, amount float64) error {
	query := "update wallet_accounts set balance = balance + ?, frozen_amount = frozen_amount - ? where user_id = ? and frozen_amount >= ?"
	result, err := m.conn.ExecCtx(ctx, query, amount, amount, userId, amount)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return zeroerr.ErrInsufficientFrozenAmount
	}
	return nil
}

func (m *customWalletAccountsModel) UpdatePayPassword(ctx context.Context, userId uint64, password string) error {
	query := "update wallet_accounts set pay_password = ? where user_id = ?"
	_, err := m.conn.ExecCtx(ctx, query, sql.NullString{String: password, Valid: true}, userId)
	return err
}

func (m *customWalletAccountsModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customWalletAccountsModel) withSession(session sqlx.Session) WalletAccountsModel {
	return &customWalletAccountsModel{
		defaultWalletAccountsModel: &defaultWalletAccountsModel{
			conn:  sqlx.NewSqlConnFromSession(session),
			table: m.table,
		},
	}
}
