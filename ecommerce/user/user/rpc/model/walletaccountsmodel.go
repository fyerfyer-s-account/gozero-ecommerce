package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ WalletAccountsModel = (*customWalletAccountsModel)(nil)

type (
	// WalletAccountsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWalletAccountsModel.
	WalletAccountsModel interface {
		walletAccountsModel
		withSession(session sqlx.Session) WalletAccountsModel
	}

	customWalletAccountsModel struct {
		*defaultWalletAccountsModel
	}
)

// NewWalletAccountsModel returns a model for the database table.
func NewWalletAccountsModel(conn sqlx.SqlConn) WalletAccountsModel {
	return &customWalletAccountsModel{
		defaultWalletAccountsModel: newWalletAccountsModel(conn),
	}
}

func (m *customWalletAccountsModel) withSession(session sqlx.Session) WalletAccountsModel {
	return NewWalletAccountsModel(sqlx.NewSqlConnFromSession(session))
}
