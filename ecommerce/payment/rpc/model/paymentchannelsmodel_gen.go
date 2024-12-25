// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

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
	paymentChannelsFieldNames          = builder.RawFieldNames(&PaymentChannels{})
	paymentChannelsRows                = strings.Join(paymentChannelsFieldNames, ",")
	paymentChannelsRowsExpectAutoSet   = strings.Join(stringx.Remove(paymentChannelsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	paymentChannelsRowsWithPlaceHolder = strings.Join(stringx.Remove(paymentChannelsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	paymentChannelsModel interface {
		Insert(ctx context.Context, data *PaymentChannels) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*PaymentChannels, error)
		FindOneByChannel(ctx context.Context, channel int64) (*PaymentChannels, error)
		Update(ctx context.Context, data *PaymentChannels) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultPaymentChannelsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	PaymentChannels struct {
		Id        uint64    `db:"id"`         // è‡ªå¢žID
		Name      string    `db:"name"`       // æ¸ é“åç§°
		Channel   int64     `db:"channel"`    // æ¸ é“ç±»åž‹ 1:å¾®ä¿¡ 2:æ”¯ä»˜å®
		Config    string    `db:"config"`     // æ¸ é“é…ç½®
		Status    int64     `db:"status"`     // çŠ¶æ€ 1:å¯ç”¨ 2:ç¦ç”¨
		CreatedAt time.Time `db:"created_at"` // åˆ›å»ºæ—¶é—´
		UpdatedAt time.Time `db:"updated_at"` // æ›´æ–°æ—¶é—´
	}
)

func newPaymentChannelsModel(conn sqlx.SqlConn) *defaultPaymentChannelsModel {
	return &defaultPaymentChannelsModel{
		conn:  conn,
		table: "`payment_channels`",
	}
}

func (m *defaultPaymentChannelsModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultPaymentChannelsModel) FindOne(ctx context.Context, id uint64) (*PaymentChannels, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", paymentChannelsRows, m.table)
	var resp PaymentChannels
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

func (m *defaultPaymentChannelsModel) FindOneByChannel(ctx context.Context, channel int64) (*PaymentChannels, error) {
	var resp PaymentChannels
	query := fmt.Sprintf("select %s from %s where `channel` = ? limit 1", paymentChannelsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, channel)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentChannelsModel) Insert(ctx context.Context, data *PaymentChannels) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, paymentChannelsRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.Channel, data.Config, data.Status)
	return ret, err
}

func (m *defaultPaymentChannelsModel) Update(ctx context.Context, newData *PaymentChannels) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, paymentChannelsRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.Name, newData.Channel, newData.Config, newData.Status, newData.Id)
	return err
}

func (m *defaultPaymentChannelsModel) tableName() string {
	return m.table
}
