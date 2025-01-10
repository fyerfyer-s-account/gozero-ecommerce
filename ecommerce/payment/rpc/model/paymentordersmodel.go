package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentOrdersModel = (*customPaymentOrdersModel)(nil)

type (
	// PaymentOrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentOrdersModel.
	PaymentOrdersModel interface {
		paymentOrdersModel
		FindOneByPaymentNo(ctx context.Context, paymentNo string) (*PaymentOrders, error)
		FindByUserId(ctx context.Context, userId uint64) ([]*PaymentOrders, error)
		FindByOrderNo(ctx context.Context, orderNo string) ([]*PaymentOrders, error)
		UpdateStatus(ctx context.Context, id uint64, status int64) error
		UpdatePartial(ctx context.Context, id uint64, updates map[string]interface{}) error
		GetPaymentsByStatus(ctx context.Context, status int64) ([]*PaymentOrders, error)
		FindByPaymentNo(ctx context.Context, paymentNo string) ([]*PaymentOrders, error)
	}

	customPaymentOrdersModel struct {
		*defaultPaymentOrdersModel
	}
)

// NewPaymentOrdersModel returns a model for the database table.
func NewPaymentOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentOrdersModel {
	return &customPaymentOrdersModel{
		defaultPaymentOrdersModel: newPaymentOrdersModel(conn, c, opts...),
	}
}

func (m *customPaymentOrdersModel) FindOneByPaymentNo(ctx context.Context, paymentNo string) (*PaymentOrders, error) {
	mallPaymentPaymentOrdersPaymentNoKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentOrdersPaymentNoPrefix, paymentNo)
	var resp PaymentOrders
	err := m.QueryRowIndexCtx(ctx, &resp, mallPaymentPaymentOrdersPaymentNoKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `payment_no` = ? limit 1", paymentOrdersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, paymentNo); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customPaymentOrdersModel) FindByUserId(ctx context.Context, userId uint64) ([]*PaymentOrders, error) {
	var resp []*PaymentOrders
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", paymentOrdersRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	return resp, err
}

func (m *customPaymentOrdersModel) FindByOrderNo(ctx context.Context, orderNo string) ([]*PaymentOrders, error) {
	var resp []*PaymentOrders
	query := fmt.Sprintf("select %s from %s where `order_no` = ?", paymentOrdersRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, orderNo)
	return resp, err
}

func (m *customPaymentOrdersModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallPaymentPaymentOrdersIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentOrdersIdPrefix, id)
	mallPaymentPaymentOrdersPaymentNoKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentOrdersPaymentNoPrefix, data.PaymentNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, mallPaymentPaymentOrdersIdKey, mallPaymentPaymentOrdersPaymentNoKey)

	return err
}

func (m *customPaymentOrdersModel) UpdatePartial(ctx context.Context, id uint64, updates map[string]interface{}) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	var sets []string
	var args []interface{}
	for k, v := range updates {
		sets = append(sets, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	args = append(args, id)

	mallPaymentPaymentOrdersIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentOrdersIdPrefix, id)
	mallPaymentPaymentOrdersPaymentNoKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentOrdersPaymentNoPrefix, data.PaymentNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(sets, ", "))
		return conn.ExecCtx(ctx, query, args...)
	}, mallPaymentPaymentOrdersIdKey, mallPaymentPaymentOrdersPaymentNoKey)

	return err
}

func (m *customPaymentOrdersModel) GetPaymentsByStatus(ctx context.Context, status int64) ([]*PaymentOrders, error) {
	var resp []*PaymentOrders
	query := fmt.Sprintf("select %s from %s where `status` = ?", paymentOrdersRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status)
	return resp, err
}

func (m *customPaymentOrdersModel) FindByPaymentNo(ctx context.Context, paymentNo string) ([]*PaymentOrders, error) {
    var resp []*PaymentOrders
    query := fmt.Sprintf("select %s from %s where `payment_no` = ?", paymentOrdersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, paymentNo)
    return resp, err
}