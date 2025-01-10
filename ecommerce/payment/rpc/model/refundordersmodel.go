package model

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefundOrdersModel = (*customRefundOrdersModel)(nil)

type (
    RefundOrdersModel interface {
        refundOrdersModel
        FindByUserId(ctx context.Context, userId uint64) ([]*RefundOrders, error)
        FindByPaymentNo(ctx context.Context, paymentNo string) ([]*RefundOrders, error)
        FindByRefundNo(ctx context.Context, refundNo string) ([]*RefundOrders, error)
        FindByOrderNo(ctx context.Context, orderNo string) ([]*RefundOrders, error)
        UpdateStatus(ctx context.Context, id uint64, status int64) error
        UpdatePartial(ctx context.Context, id uint64, updates map[string]interface{}) error
        GetRefundsByStatus(ctx context.Context, status int64) ([]*RefundOrders, error)
    }

    customRefundOrdersModel struct {
        *defaultRefundOrdersModel
    }
)

func NewRefundOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RefundOrdersModel {
    return &customRefundOrdersModel{
        defaultRefundOrdersModel: newRefundOrdersModel(conn, c, opts...),
    }
}

func (m *customRefundOrdersModel) FindByUserId(ctx context.Context, userId uint64) ([]*RefundOrders, error) {
    var resp []*RefundOrders
    query := fmt.Sprintf("select %s from %s where `user_id` = ?", refundOrdersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
    return resp, err
}

func (m *customRefundOrdersModel) FindByPaymentNo(ctx context.Context, paymentNo string) ([]*RefundOrders, error) {
    var resp []*RefundOrders
    query := fmt.Sprintf("select %s from %s where `payment_no` = ?", refundOrdersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, paymentNo)
    return resp, err
}

func (m *customRefundOrdersModel) FindByOrderNo(ctx context.Context, orderNo string) ([]*RefundOrders, error) {
    var resp []*RefundOrders
    query := fmt.Sprintf("select %s from %s where `order_no` = ?", refundOrdersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, orderNo)
    return resp, err
}

func (m *customRefundOrdersModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
    data, err := m.FindOne(ctx, id)
    if err != nil {
        return err
    }

    mallPaymentRefundOrdersIdKey := fmt.Sprintf("%s%v", cacheMallPaymentRefundOrdersIdPrefix, id)
    mallPaymentRefundOrdersRefundNoKey := fmt.Sprintf("%s%v", cacheMallPaymentRefundOrdersRefundNoPrefix, data.RefundNo)
    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, id)
    }, mallPaymentRefundOrdersIdKey, mallPaymentRefundOrdersRefundNoKey)

    return err
}

func (m *customRefundOrdersModel) UpdatePartial(ctx context.Context, id uint64, updates map[string]interface{}) error {
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

    mallPaymentRefundOrdersIdKey := fmt.Sprintf("%s%v", cacheMallPaymentRefundOrdersIdPrefix, id)
    mallPaymentRefundOrdersRefundNoKey := fmt.Sprintf("%s%v", cacheMallPaymentRefundOrdersRefundNoPrefix, data.RefundNo)
    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(sets, ", "))
        return conn.ExecCtx(ctx, query, args...)
    }, mallPaymentRefundOrdersIdKey, mallPaymentRefundOrdersRefundNoKey)

    return err
}

func (m *customRefundOrdersModel) GetRefundsByStatus(ctx context.Context, status int64) ([]*RefundOrders, error) {
    var resp []*RefundOrders
    query := fmt.Sprintf("select %s from %s where `status` = ?", refundOrdersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status)
    return resp, err
}

func (m *customRefundOrdersModel) FindByRefundNo(ctx context.Context, refundNo string) ([]*RefundOrders, error) {
    var resp []*RefundOrders
    query := fmt.Sprintf("select %s from %s where `refund_no` = ?", refundOrdersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, refundNo)
    return resp, err
}