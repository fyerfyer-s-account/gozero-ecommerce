package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WalletTransactionsModel = (*customWalletTransactionsModel)(nil)

type (
	WalletTransactionsModel interface {
		walletTransactionsModel
		WithSession(session sqlx.Session) WalletTransactionsModel
		FindByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*WalletTransactions, error)
		FindByType(ctx context.Context, userId uint64, transType int64) ([]*WalletTransactions, error)
		GetTransactionStats(ctx context.Context, tranId uint64) (*TransactionStats, error)
		UpdateState(ctx context.Context, status int64, tranId uint64) (int64, error)
		DeleteByUserId(ctx context.Context, userId uint64) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
	}

	customWalletTransactionsModel struct {
		*defaultWalletTransactionsModel
	}

	TransactionStats struct {
		TotalRecharge float64 `db:"total_recharge"`
		TotalSpent    float64 `db:"total_spent"`
		TotalRefund   float64 `db:"total_refund"`
	}
)

func NewWalletTransactionsModel(conn sqlx.SqlConn) WalletTransactionsModel {
	return &customWalletTransactionsModel{
		defaultWalletTransactionsModel: newWalletTransactionsModel(conn),
	}
}

func (m *customWalletTransactionsModel) FindByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*WalletTransactions, error) {
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by created_at desc limit ?,?",
		walletTransactionsRows, m.table)
	var resp []*WalletTransactions
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customWalletTransactionsModel) FindByType(ctx context.Context, userId uint64, transType int64) ([]*WalletTransactions, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `type` = ? order by created_at desc",
		walletTransactionsRows, m.table)
	var resp []*WalletTransactions
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, transType)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customWalletTransactionsModel) UpdateState(ctx context.Context, status int64, tranId uint64) (int64, error) {
	query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
	res, err := m.conn.ExecCtx(ctx, query, status, tranId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *customWalletTransactionsModel) GetTransactionStats(ctx context.Context, userId uint64) (*TransactionStats, error) {
	query := fmt.Sprintf(`
        SELECT 
            COALESCE(SUM(CASE WHEN type = 1 THEN amount ELSE 0 END), 0) as total_recharge,
            COALESCE(SUM(CASE WHEN type = 3 THEN amount ELSE 0 END), 0) as total_spent,
            COALESCE(SUM(CASE WHEN type = 4 THEN amount ELSE 0 END), 0) as total_refund
        FROM %s 
        WHERE user_id = ? AND status = 1`, m.table)

	var stats TransactionStats
	err := m.conn.QueryRowCtx(ctx, &stats, query, userId)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (m *customWalletTransactionsModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customWalletTransactionsModel) WithSession(session sqlx.Session) WalletTransactionsModel {
	return &customWalletTransactionsModel{
		defaultWalletTransactionsModel: &defaultWalletTransactionsModel{
			conn:  sqlx.NewSqlConnFromSession(session),
			table: m.table,
		},
	}
}

func (m *customWalletTransactionsModel) DeleteByUserId(ctx context.Context, userId uint64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}
