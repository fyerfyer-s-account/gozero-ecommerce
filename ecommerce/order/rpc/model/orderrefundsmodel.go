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

var _ OrderRefundsModel = (*customOrderRefundsModel)(nil)

type (
	// OrderRefundsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderRefundsModel.
	OrderRefundsModel interface {
		orderRefundsModel
		FindByOrderId(ctx context.Context, orderId uint64) ([]*OrderRefunds, error)
        UpdateStatus(ctx context.Context, refundNo string, status int64, reply string) error
        GetLatestRefund(ctx context.Context, orderId uint64) (*OrderRefunds, error)
        FindByStatus(ctx context.Context, status int64) ([]*OrderRefunds, error)
	}

	customOrderRefundsModel struct {
		*defaultOrderRefundsModel
	}
)

// NewOrderRefundsModel returns a model for the database table.
func NewOrderRefundsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderRefundsModel {
	return &customOrderRefundsModel{
		defaultOrderRefundsModel: newOrderRefundsModel(conn, c, opts...),
	}
}

func (m *customOrderRefundsModel) FindByOrderId(ctx context.Context, orderId uint64) ([]*OrderRefunds, error) {
    var resp []*OrderRefunds
    query := fmt.Sprintf("select %s from %s where `order_id` = ?", orderRefundsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, orderId)
    return resp, err
}

func (m *customOrderRefundsModel) UpdateStatus(ctx context.Context, refundNo string, status int64, reply string) error {
    refund, err := m.FindOneByRefundNo(ctx, refundNo)
    if err != nil {
        return err
    }

    refundIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsIdPrefix, refund.Id)
    refundNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsRefundNoPrefix, refundNo)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ?, `reply` = ?, `updated_at` = ? where `refund_no` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, reply, time.Now(), refundNo)
    }, refundIdKey, refundNoKey)

    return err
}

func (m *customOrderRefundsModel) GetLatestRefund(ctx context.Context, orderId uint64) (*OrderRefunds, error) {
    var resp OrderRefunds
    query := fmt.Sprintf("select %s from %s where `order_id` = ? order by created_at desc limit 1", orderRefundsRows, m.table)
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

func (m *customOrderRefundsModel) FindByStatus(ctx context.Context, status int64) ([]*OrderRefunds, error) {
    var resp []*OrderRefunds
    query := fmt.Sprintf("select %s from %s where `status` = ?", orderRefundsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status)
    return resp, err
}