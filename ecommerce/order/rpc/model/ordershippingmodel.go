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

var _ OrderShippingModel = (*customOrderShippingModel)(nil)

type (
	// OrderShippingModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderShippingModel.
	OrderShippingModel interface {
		orderShippingModel
        FindByStatus(ctx context.Context, status int64) ([]*OrderShipping, error)
        FindByOrderId(ctx context.Context, orderId uint64) (*OrderShipping, error)
        UpdateStatus(ctx context.Context, id uint64, status int64) error
        UpdateShippingInfo(ctx context.Context, orderId uint64, shippingNo, company string) error
	}

	customOrderShippingModel struct {
		*defaultOrderShippingModel
	}
)

// NewOrderShippingModel returns a model for the database table.
func NewOrderShippingModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderShippingModel {
	return &customOrderShippingModel{
		defaultOrderShippingModel: newOrderShippingModel(conn, c, opts...),
	}
}

func (m *customOrderShippingModel) FindByOrderId(ctx context.Context, orderId uint64) (*OrderShipping, error) {
    var resp OrderShipping
    query := fmt.Sprintf("select %s from %s where `order_id` = ? limit 1", orderShippingRows, m.table)
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

func (m *customOrderShippingModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
    data, err := m.FindOne(ctx, id)
    if err != nil {
        return err
    }

    shippingIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderShippingIdPrefix, id)
    shippingNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderShippingShippingNoPrefix, data.ShippingNo)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ?, `updated_at` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, time.Now(), id)
    }, shippingIdKey, shippingNoKey)

    return err
}

func (m *customOrderShippingModel) UpdateShippingInfo(ctx context.Context, orderId uint64, shippingNo, company string) error {
    shipping, err := m.FindByOrderId(ctx, orderId)
    if err != nil {
        return err
    }

    shippingIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderShippingIdPrefix, shipping.Id)
    oldShippingNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderShippingShippingNoPrefix, shipping.ShippingNo)
    newShippingNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderShippingShippingNoPrefix, shippingNo)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `shipping_no` = ?, `company` = ?, `status` = ?, `ship_time` = ?, `updated_at` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, shippingNo, company, 1, time.Now(), time.Now(), shipping.Id)
    }, shippingIdKey, oldShippingNoKey, newShippingNoKey)

    return err
}

func (m *customOrderShippingModel) FindByStatus(ctx context.Context, status int64) ([]*OrderShipping, error) {
    var resp []*OrderShipping
    query := fmt.Sprintf("select %s from %s where `status` = ?", orderShippingRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status)
    return resp, err
}