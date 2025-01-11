package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrdersModel = (*customOrdersModel)(nil)

type (
    OrdersModel interface {
        ordersModel
        FindByUserIdAndStatus(ctx context.Context, userId uint64, status int64) ([]*Orders, error)
        FindByStatus(ctx context.Context, status int64) ([]*Orders, error)
        UpdateStatus(ctx context.Context, id uint64, status int64) error
        FindByUserIdWithPage(ctx context.Context, userId uint64, status int64, page, pageSize int) ([]*Orders, error)
        CountByUserIdAndStatus(ctx context.Context, userId uint64, status int64) (int64, error)
        BatchUpdateStatus(ctx context.Context, ids []uint64, status int64) error
        CreateOrder(ctx context.Context, order *Orders) (uint64, error)
    }

    customOrdersModel struct {
        *defaultOrdersModel
    }
)

func NewOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrdersModel {
    return &customOrdersModel{
        defaultOrdersModel: newOrdersModel(conn, c, opts...),
    }
}

func (m *customOrdersModel) FindByUserIdAndStatus(ctx context.Context, userId uint64, status int64) ([]*Orders, error) {
    var resp []*Orders
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and `status` = ?", ordersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, status)
    return resp, err
}

func (m *customOrdersModel) FindByStatus(ctx context.Context, status int64) ([]*Orders, error) {
    var resp []*Orders
    query := fmt.Sprintf("select %s from %s where `status` = ?", ordersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, status)
    return resp, err
}

func (m *customOrdersModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
    data, err := m.FindOne(ctx, id)
    if err != nil {
        return err
    }

    orderIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersIdPrefix, id)
    orderNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersOrderNoPrefix, data.OrderNo)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ?, `updated_at` = now() where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, id)
    }, orderIdKey, orderNoKey)

    return err
}

func (m *customOrdersModel) FindByUserIdWithPage(ctx context.Context, userId uint64, status int64, page, pageSize int) ([]*Orders, error) {
    if page < 1 {
        page = 1
    }
    offset := (page - 1) * pageSize

    var resp []*Orders
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and (`status` = ? or ? = -1) order by `created_at` desc limit ?, ?", 
        ordersRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, status, status, offset, pageSize)
    return resp, err
}

func (m *customOrdersModel) CountByUserIdAndStatus(ctx context.Context, userId uint64, status int64) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and (`status` = ? or ? = -1)", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, userId, status, status)
    return count, err
}

func (m *customOrdersModel) BatchUpdateStatus(ctx context.Context, ids []uint64, status int64) error {
    if len(ids) == 0 {
        return nil
    }

    args := make([]interface{}, 0, len(ids)+1)
    args = append(args, status)
    placeholder := make([]string, len(ids))
    for i, id := range ids {
        placeholder[i] = "?"
        args = append(args, id)
    }

    query := fmt.Sprintf("update %s set `status` = ?, `updated_at` = now() where `id` in (%s)", 
        m.table, strings.Join(placeholder, ","))
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })

    return err
}

func (m *customOrdersModel) CreateOrder(ctx context.Context, order *Orders) (uint64, error) {
    result, err := m.Insert(ctx, order)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    return uint64(id), err
}