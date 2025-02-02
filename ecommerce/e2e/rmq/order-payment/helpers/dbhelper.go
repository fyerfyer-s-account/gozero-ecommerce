package helpers

import (
    "context"
    "fmt"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DBHelper struct {
    conn          sqlx.SqlConn
    cacheRedis    cache.ClusterConf
    ordersModel   model.OrdersModel
    paymentsModel model.OrderPaymentsModel
    refundsModel  model.OrderRefundsModel
}

// NewDBHelper creates a new database helper with the given configuration
func NewDBHelper(mysqlDSN string, redisConf cache.ClusterConf) (*DBHelper, error) {
    conn := sqlx.NewMysql(mysqlDSN)

    return &DBHelper{
        conn:          conn,
        cacheRedis:    redisConf,
        ordersModel:   model.NewOrdersModel(conn, redisConf),
        paymentsModel: model.NewOrderPaymentsModel(conn, redisConf),
        refundsModel:  model.NewOrderRefundsModel(conn, redisConf),
    }, nil
}

// Close closes the database connections
func (h *DBHelper) Close() error {
    return nil // sqlx doesn't provide a Close method
}

// GetOrdersModel returns the orders model
func (h *DBHelper) GetOrdersModel() model.OrdersModel {
    return h.ordersModel
}

// GetPaymentsModel returns the payments model 
func (h *DBHelper) GetPaymentsModel() model.OrderPaymentsModel {
    return h.paymentsModel
}

// GetRefundsModel returns the refunds model
func (h *DBHelper) GetRefundsModel() model.OrderRefundsModel {
    return h.refundsModel
}

// CleanTestData removes test data from all relevant tables
func (h *DBHelper) CleanTestData(ctx context.Context) error {
    if _, err := h.conn.ExecCtx(ctx, "DELETE FROM orders WHERE order_no LIKE 'TEST_%'"); err != nil {
        return fmt.Errorf("failed to clean table orders: %w", err)
    }

    if _, err := h.conn.ExecCtx(ctx, "DELETE FROM order_payments WHERE payment_no LIKE 'TEST_%'"); err != nil {
        return fmt.Errorf("failed to clean table order_payments: %w", err)
    }

    if _, err := h.conn.ExecCtx(ctx, "DELETE FROM order_refunds WHERE refund_no LIKE 'TEST_%'"); err != nil {
        return fmt.Errorf("failed to clean table order_refunds: %w", err)
    }

    return nil
}