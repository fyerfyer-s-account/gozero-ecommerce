package model

import (
	"context"
    "database/sql"
    "fmt"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CartStatisticsModel = (*customCartStatisticsModel)(nil)

type (
	// CartStatisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartStatisticsModel.
	CartStatisticsModel interface {
		cartStatisticsModel
        UpdateTotalQuantity(ctx context.Context, userId uint64, quantity int64) error
        UpdateTotalAmount(ctx context.Context, userId uint64, amount float64) error
        UpdateSelectedStats(ctx context.Context, userId uint64, quantity int64, amount float64) error
        RecalculateStats(ctx context.Context, userId uint64) error
        Upsert(ctx context.Context, data *CartStatistics) error
	}

	customCartStatisticsModel struct {
		*defaultCartStatisticsModel
	}
)

// NewCartStatisticsModel returns a model for the database table.
func NewCartStatisticsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CartStatisticsModel {
	return &customCartStatisticsModel{
		defaultCartStatisticsModel: newCartStatisticsModel(conn, c, opts...),
	}
}

func (m *customCartStatisticsModel) UpdateTotalQuantity(ctx context.Context, userId uint64, quantity int64) error {
    key := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, userId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, 
            fmt.Sprintf("update %s set total_quantity = total_quantity + ?, updated_at = now() where user_id = ?", 
            m.table), quantity, userId)
    }, key)
    return err
}

func (m *customCartStatisticsModel) UpdateTotalAmount(ctx context.Context, userId uint64, amount float64) error {
    key := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, userId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, 
            fmt.Sprintf("update %s set total_amount = total_amount + ?, updated_at = now() where user_id = ?", 
            m.table), amount, userId)
    }, key)
    return err
}

func (m *customCartStatisticsModel) UpdateSelectedStats(ctx context.Context, userId uint64, quantity int64, amount float64) error {
    key := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, userId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, 
            fmt.Sprintf("update %s set selected_quantity = selected_quantity + ?, selected_amount = selected_amount + ?, updated_at = now() where user_id = ?", 
            m.table), quantity, amount, userId)
    }, key)
    return err
}

func (m *customCartStatisticsModel) RecalculateStats(ctx context.Context, userId uint64) error {
    key := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, userId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, `
            update cart_statistics cs
            set cs.total_quantity = (
                select IFNULL(sum(quantity), 0) from cart_items where user_id = ?
            ),
            cs.total_amount = (
                select IFNULL(sum(price * quantity), 0) from cart_items where user_id = ?
            ),
            cs.selected_quantity = (
                select IFNULL(sum(quantity), 0) from cart_items where user_id = ? and selected = 1
            ),
            cs.selected_amount = (
                select IFNULL(sum(price * quantity), 0) from cart_items where user_id = ? and selected = 1
            ),
            cs.updated_at = now()
            where cs.user_id = ?
        `, userId, userId, userId, userId, userId)
    }, key)
    return err
}

func (m *customCartStatisticsModel) Upsert(ctx context.Context, data *CartStatistics) error {
    key := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, data.UserId)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, `
            insert into cart_statistics (
                user_id, total_quantity, selected_quantity, 
                total_amount, selected_amount, updated_at
            ) values (?, ?, ?, ?, ?, now())
            on duplicate key update
                total_quantity = values(total_quantity),
                selected_quantity = values(selected_quantity),
                total_amount = values(total_amount),
                selected_amount = values(selected_amount),
                updated_at = now()
        `, data.UserId, data.TotalQuantity, data.SelectedQuantity,
            data.TotalAmount, data.SelectedAmount)
    }, key)
    return err
}