package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WalletTransactionsModel = (*customWalletTransactionsModel)(nil)

type (
	// WalletTransactionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWalletTransactionsModel.
	WalletTransactionsModel interface {
		walletTransactionsModel
		withSession(session sqlx.Session) WalletTransactionsModel
	}

	customWalletTransactionsModel struct {
		*defaultWalletTransactionsModel
	}
)

// NewWalletTransactionsModel returns a model for the database table.
func NewWalletTransactionsModel(conn sqlx.SqlConn) WalletTransactionsModel {
	return &customWalletTransactionsModel{
		defaultWalletTransactionsModel: newWalletTransactionsModel(conn),
	}
}

func (m *customWalletTransactionsModel) withSession(session sqlx.Session) WalletTransactionsModel {
	return NewWalletTransactionsModel(sqlx.NewSqlConnFromSession(session))
}
