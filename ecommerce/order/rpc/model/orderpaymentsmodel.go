package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderPaymentsModel = (*customOrderPaymentsModel)(nil)

type (
	// OrderPaymentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderPaymentsModel.
	OrderPaymentsModel interface {
		orderPaymentsModel
		FindByOrderId(ctx context.Context, orderId uint64) (*OrderPayments, error)
        UpdateStatus(ctx context.Context, paymentNo string, status int64, payTime time.Time) error
        FindByStatusAndTime(ctx context.Context, status int64, startTime, endTime time.Time) ([]*OrderPayments, error)
        CreatePayment(ctx context.Context, payment *OrderPayments) (uint64, error)
	}

	customOrderPaymentsModel struct {
		*defaultOrderPaymentsModel
	}
)

// NewOrderPaymentsModel returns a model for the database table.
func NewOrderPaymentsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderPaymentsModel {
	return &customOrderPaymentsModel{
		defaultOrderPaymentsModel: newOrderPaymentsModel(conn, c, opts...),
	}
}

func (m *customOrderPaymentsModel) FindByOrderId(ctx context.Context, orderId uint64) (*OrderPayments, error) {
    var resp OrderPayments
    query := fmt.Sprintf("select %s from %s where `order_id` = ? limit 1", orderPaymentsRows, m.table)
    err := m.QueryRowNoCacheCtx(ctx, &resp, query, orderId)
    switch err {
    case nil:
        return &resp, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *customOrderPaymentsModel) UpdateStatus(ctx context.Context, paymentNo string, status int64, payTime time.Time) error {
    payment, err := m.FindOneByPaymentNo(ctx, paymentNo)
    if err != nil {
        return err
    }

    paymentIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderPaymentsIdPrefix, payment.Id)
    paymentNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderPaymentsPaymentNoPrefix, paymentNo)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ?, `pay_time` = ?, `updated_at` = ? where `payment_no` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, payTime, time.Now(), paymentNo)
    }, paymentIdKey, paymentNoKey)

    return err
}

func (m *customOrderPaymentsModel) FindByStatusAndTime(ctx context.Context, status int64, startTime, endTime time.Time) ([]*OrderPayments, error) {
    var resp []*OrderPayments
    query := fmt.Sprintf("select %s from %s where `status` = ? and `created_at` between ? and ?", orderPaymentsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status, startTime, endTime)
    return resp, err
}

func (m *customOrderPaymentsModel) CreatePayment(ctx context.Context, payment *OrderPayments) (uint64, error) {
    result, err := m.Insert(ctx, payment)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    return uint64(id), err
}