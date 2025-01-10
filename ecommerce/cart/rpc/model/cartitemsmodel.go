package model

import (
	"context"
    "database/sql"
    "fmt"
    "strings"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CartItemsModel = (*customCartItemsModel)(nil)

type (
	// CartItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartItemsModel.
	CartItemsModel interface {
		cartItemsModel
		FindByUserId(ctx context.Context, userId uint64) ([]*CartItems, error)
        FindSelectedByUserId(ctx context.Context, userId uint64) ([]*CartItems, error)
        UpdateQuantity(ctx context.Context, userId, productId, skuId uint64, quantity int64) error
        UpdateSelected(ctx context.Context, userId, productId, skuId uint64, selected int64) error
        UpdateAllSelected(ctx context.Context, userId uint64, selected int64) error
        DeleteByUserId(ctx context.Context, userId uint64) error
        BatchInsert(ctx context.Context, items []*CartItems) error
	}

	customCartItemsModel struct {
		*defaultCartItemsModel
	}
)

// NewCartItemsModel returns a model for the database table.
func NewCartItemsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CartItemsModel {
	return &customCartItemsModel{
		defaultCartItemsModel: newCartItemsModel(conn, c, opts...),
	}
}

func (m *customCartItemsModel) FindByUserId(ctx context.Context, userId uint64) ([]*CartItems, error) {
    var resp []*CartItems
    query := fmt.Sprintf("select %s from %s where `user_id` = ?", cartItemsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
    return resp, err
}

func (m *customCartItemsModel) FindSelectedByUserId(ctx context.Context, userId uint64) ([]*CartItems, error) {
    var resp []*CartItems
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and `selected` = 1", cartItemsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
    return resp, err
}

func (m *customCartItemsModel) UpdateQuantity(ctx context.Context, userId, productId, skuId uint64, quantity int64) error {
    cartItem, err := m.FindOneByUserIdSkuId(ctx, userId, skuId)
    if err != nil {
        return err
    }

    cartItemKey := fmt.Sprintf("%s%v", cacheMallCartCartItemsIdPrefix, cartItem.Id)
    cartItemUserSkuKey := fmt.Sprintf("%s%v:%v", cacheMallCartCartItemsUserIdSkuIdPrefix, userId, skuId)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `quantity` = ?, `updated_at` = now() where `user_id` = ? and `product_id` = ? and `sku_id` = ?", m.table)
        return conn.ExecCtx(ctx, query, quantity, userId, productId, skuId)
    }, cartItemKey, cartItemUserSkuKey)

    return err
}

func (m *customCartItemsModel) UpdateSelected(ctx context.Context, userId, productId, skuId uint64, selected int64) error {
    cartItem, err := m.FindOneByUserIdSkuId(ctx, userId, skuId)
    if err != nil {
        return err
    }

    cartItemKey := fmt.Sprintf("%s%v", cacheMallCartCartItemsIdPrefix, cartItem.Id)
    cartItemUserSkuKey := fmt.Sprintf("%s%v:%v", cacheMallCartCartItemsUserIdSkuIdPrefix, userId, skuId)

    _, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `selected` = ?, `updated_at` = now() where `user_id` = ? and `product_id` = ? and `sku_id` = ?", m.table)
        return conn.ExecCtx(ctx, query, selected, userId, productId, skuId)
    }, cartItemKey, cartItemUserSkuKey)

    return err
}

func (m *customCartItemsModel) UpdateAllSelected(ctx context.Context, userId uint64, selected int64) error {
    query := fmt.Sprintf("update %s set `selected` = ?, `updated_at` = now() where `user_id` = ?", m.table)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, selected, userId)
    })
    return err
}

func (m *customCartItemsModel) DeleteByUserId(ctx context.Context, userId uint64) error {
    query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, userId)
    })
    return err
}

func (m *customCartItemsModel) BatchInsert(ctx context.Context, items []*CartItems) error {
    if len(items) == 0 {
        return nil
    }

    values := make([]string, 0, len(items))
    args := make([]interface{}, 0, len(items)*9)
    for _, item := range items {
        values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
        args = append(args, item.UserId, item.ProductId, item.SkuId, item.ProductName, 
            item.SkuName, item.Image, item.Price, item.Quantity, item.Selected)
    }

    query := fmt.Sprintf("insert into %s (%s) values %s", 
        m.table, cartItemsRowsExpectAutoSet, strings.Join(values, ","))
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })

    return err
}