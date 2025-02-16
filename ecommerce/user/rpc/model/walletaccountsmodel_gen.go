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
	walletAccountsFieldNames          = builder.RawFieldNames(&WalletAccounts{})
	walletAccountsRows                = strings.Join(walletAccountsFieldNames, ",")
	walletAccountsRowsExpectAutoSet   = strings.Join(stringx.Remove(walletAccountsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	walletAccountsRowsWithPlaceHolder = strings.Join(stringx.Remove(walletAccountsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	walletAccountsModel interface {
		Insert(ctx context.Context, data *WalletAccounts) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*WalletAccounts, error)
		FindOneByUserId(ctx context.Context, userId uint64) (*WalletAccounts, error)
		Update(ctx context.Context, data *WalletAccounts) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultWalletAccountsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	WalletAccounts struct {
		Id           uint64         `db:"id"`            // é’±åŒ…ID
		UserId       uint64         `db:"user_id"`       // ç”¨æˆ·ID
		Balance      float64        `db:"balance"`       // è´¦æˆ·ä½™é¢
		FrozenAmount float64        `db:"frozen_amount"` // å†»ç»“é‡‘é¢
		PayPassword  sql.NullString `db:"pay_password"`  // æ”¯ä»˜å¯†ç 
		Status       int64          `db:"status"`        // çŠ¶æ€ 0:å†»ç»“ 1:æ­£å¸¸
		CreatedAt    time.Time      `db:"created_at"`    // åˆ›å»ºæ—¶é—´
		UpdatedAt    time.Time      `db:"updated_at"`    // æ›´æ–°æ—¶é—´
	}
)

func newWalletAccountsModel(conn sqlx.SqlConn) *defaultWalletAccountsModel {
	return &defaultWalletAccountsModel{
		conn:  conn,
		table: "`wallet_accounts`",
	}
}

func (m *defaultWalletAccountsModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultWalletAccountsModel) FindOne(ctx context.Context, id uint64) (*WalletAccounts, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", walletAccountsRows, m.table)
	var resp WalletAccounts
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

func (m *defaultWalletAccountsModel) FindOneByUserId(ctx context.Context, userId uint64) (*WalletAccounts, error) {
	var resp WalletAccounts
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", walletAccountsRows, m.table)
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

func (m *defaultWalletAccountsModel) Insert(ctx context.Context, data *WalletAccounts) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, walletAccountsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Balance, data.FrozenAmount, data.PayPassword, data.Status)
	return ret, err
}

func (m *defaultWalletAccountsModel) Update(ctx context.Context, newData *WalletAccounts) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, walletAccountsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.UserId, newData.Balance, newData.FrozenAmount, newData.PayPassword, newData.Status, newData.Id)
	return err
}

func (m *defaultWalletAccountsModel) tableName() string {
	return m.table
}
