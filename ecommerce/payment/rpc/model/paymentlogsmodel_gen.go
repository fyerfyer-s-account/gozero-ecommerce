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
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	paymentLogsFieldNames          = builder.RawFieldNames(&PaymentLogs{})
	paymentLogsRows                = strings.Join(paymentLogsFieldNames, ",")
	paymentLogsRowsExpectAutoSet   = strings.Join(stringx.Remove(paymentLogsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	paymentLogsRowsWithPlaceHolder = strings.Join(stringx.Remove(paymentLogsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallPaymentPaymentLogsIdPrefix = "cache:mallPayment:paymentLogs:id:"
)

type (
	paymentLogsModel interface {
		Insert(ctx context.Context, data *PaymentLogs) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*PaymentLogs, error)
		Update(ctx context.Context, data *PaymentLogs) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultPaymentLogsModel struct {
		sqlc.CachedConn
		table string
	}

	PaymentLogs struct {
		Id           uint64         `db:"id"`            // è‡ªå¢žID
		PaymentNo    string         `db:"payment_no"`    // æ”¯ä»˜å•å·
		Type         int64          `db:"type"`          // ç±»åž‹ 1:æ”¯ä»˜ 2:é€€æ¬¾
		Channel      int64          `db:"channel"`       // æ”¯ä»˜æ¸ é“
		RequestData  sql.NullString `db:"request_data"`  // è¯·æ±‚æ•°æ®
		ResponseData sql.NullString `db:"response_data"` // å“åº”æ•°æ®
		CreatedAt    time.Time      `db:"created_at"`    // åˆ›å»ºæ—¶é—´
	}
)

func newPaymentLogsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultPaymentLogsModel {
	return &defaultPaymentLogsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`payment_logs`",
	}
}

func (m *defaultPaymentLogsModel) Delete(ctx context.Context, id uint64) error {
	mallPaymentPaymentLogsIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentLogsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mallPaymentPaymentLogsIdKey)
	return err
}

func (m *defaultPaymentLogsModel) FindOne(ctx context.Context, id uint64) (*PaymentLogs, error) {
	mallPaymentPaymentLogsIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentLogsIdPrefix, id)
	var resp PaymentLogs
	err := m.QueryRowCtx(ctx, &resp, mallPaymentPaymentLogsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", paymentLogsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPaymentLogsModel) Insert(ctx context.Context, data *PaymentLogs) (sql.Result, error) {
	mallPaymentPaymentLogsIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentLogsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, paymentLogsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.PaymentNo, data.Type, data.Channel, data.RequestData, data.ResponseData)
	}, mallPaymentPaymentLogsIdKey)
	return ret, err
}

func (m *defaultPaymentLogsModel) Update(ctx context.Context, data *PaymentLogs) error {
	mallPaymentPaymentLogsIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentLogsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, paymentLogsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.PaymentNo, data.Type, data.Channel, data.RequestData, data.ResponseData, data.Id)
	}, mallPaymentPaymentLogsIdKey)
	return err
}

func (m *defaultPaymentLogsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallPaymentPaymentLogsIdPrefix, primary)
}

func (m *defaultPaymentLogsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", paymentLogsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultPaymentLogsModel) tableName() string {
	return m.table
}
