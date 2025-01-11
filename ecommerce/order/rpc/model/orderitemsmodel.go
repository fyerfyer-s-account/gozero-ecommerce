package model

import (
	"context"
    "database/sql"
    "fmt"
    "strings"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemsModel = (*customOrderItemsModel)(nil)

type (
	// OrderItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemsModel.
	OrderItemsModel interface {
		orderItemsModel
		FindByOrderId(ctx context.Context, orderId uint64) ([]*OrderItems, error)
        BatchInsert(ctx context.Context, items []*OrderItems) error
        CalculateOrderTotal(ctx context.Context, orderId uint64) (float64, error)
	}

	customOrderItemsModel struct {
		*defaultOrderItemsModel
	}
)

// NewOrderItemsModel returns a model for the database table.
func NewOrderItemsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderItemsModel {
	return &customOrderItemsModel{
		defaultOrderItemsModel: newOrderItemsModel(conn, c, opts...),
	}
}

func (m *customOrderItemsModel) FindByOrderId(ctx context.Context, orderId uint64) ([]*OrderItems, error) {
    var resp []*OrderItems
    query := fmt.Sprintf("select %s from %s where `order_id` = ?", orderItemsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, orderId)
    return resp, err
}

func (m *customOrderItemsModel) BatchInsert(ctx context.Context, items []*OrderItems) error {
    if len(items) == 0 {
        return nil
    }

    values := make([]string, 0, len(items))
    args := make([]interface{}, 0, len(items)*8)
    for _, item := range items {
        values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?)")
        args = append(args, item.OrderId, item.ProductId, item.SkuId, 
            item.ProductName, item.SkuName, item.Price, item.Quantity, item.TotalAmount)
    }

    query := fmt.Sprintf("insert into %s (%s) values %s", 
        m.table, orderItemsRowsExpectAutoSet, strings.Join(values, ","))
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })

    return err
}

func (m *customOrderItemsModel) CalculateOrderTotal(ctx context.Context, orderId uint64) (float64, error) {
    var total float64
    query := fmt.Sprintf("select sum(total_amount) from %s where `order_id` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &total, query, orderId)
    return total, err
}