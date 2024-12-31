package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ SkusModel = (*customSkusModel)(nil)

type (
	SkusModel interface {
		skusModel
		FindManyByProductId(ctx context.Context, productId uint64) ([]*Skus, error)
		FindManyPageByProductId(ctx context.Context, productId uint64, page, pageSize int) ([]*Skus, error)
		Count(ctx context.Context, productId uint64) (int64, error)
		UpdateStock(ctx context.Context, id uint64, increment int64) error
		UpdateSales(ctx context.Context, id uint64, increment int64) error
		UpdatePriceAndStock(ctx context.Context, id uint64, price float64, stock int64) error
		BatchInsert(ctx context.Context, data []*Skus) error
		DeleteByProductId(ctx context.Context, productId uint64) error
	}

	customSkusModel struct {
		*defaultSkusModel
	}
)

func NewSkusModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SkusModel {
	return &customSkusModel{
		defaultSkusModel: newSkusModel(conn, c, opts...),
	}
}

func (m *customSkusModel) FindManyByProductId(ctx context.Context, productId uint64) ([]*Skus, error) {
	var skus []*Skus
	query := fmt.Sprintf("select %s from %s where `product_id` = ?", skusRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &skus, query, productId)
	return skus, err
}

func (m *customSkusModel) FindManyPageByProductId(ctx context.Context, productId uint64, page, pageSize int) ([]*Skus, error) {
	var skus []*Skus
	query := fmt.Sprintf("select %s from %s where `product_id` = ? limit ?, ?", skusRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &skus, query, productId, (page-1)*pageSize, pageSize)
	return skus, err
}

func (m *customSkusModel) Count(ctx context.Context, productId uint64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `product_id` = ?", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, productId)
	return count, err
}

func (m *customSkusModel) UpdateStock(ctx context.Context, id uint64, increment int64) error {
	mallProductSkusIdKey := fmt.Sprintf("%s%v", cacheMallProductSkusIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `stock` = `stock` + ? where `id` = ? and `stock` + ? >= 0", m.table)
		return conn.ExecCtx(ctx, query, increment, id, increment)
	}, mallProductSkusIdKey)
	return err
}

func (m *customSkusModel) UpdateSales(ctx context.Context, id uint64, increment int64) error {
	mallProductSkusIdKey := fmt.Sprintf("%s%v", cacheMallProductSkusIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `sales` = `sales` + ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, increment, id)
	}, mallProductSkusIdKey)
	return err
}

func (m *customSkusModel) UpdatePriceAndStock(ctx context.Context, id uint64, price float64, stock int64) error {
	sku, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallProductSkusIdKey := fmt.Sprintf("%s%v", cacheMallProductSkusIdPrefix, id)
	mallProductSkusSkuCodeKey := fmt.Sprintf("%s%v", cacheMallProductSkusSkuCodePrefix, sku.SkuCode)

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `price` = ?, `stock` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, price, stock, id)
	}, mallProductSkusIdKey, mallProductSkusSkuCodeKey)
	return err
}

func (m *customSkusModel) BatchInsert(ctx context.Context, data []*Skus) error {
	if len(data) == 0 {
		return nil
	}

	values := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)*6)
	for _, sku := range data {
		values = append(values, "(?, ?, ?, ?, ?, ?)")
		args = append(args, sku.ProductId, sku.SkuCode, sku.Attributes, sku.Price, sku.Stock, sku.Sales)
	}

	query := fmt.Sprintf("insert into %s (%s) values %s",
		m.table, skusRowsExpectAutoSet, strings.Join(values, ","))

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	})
	return err
}

func (m *customSkusModel) DeleteByProductId(ctx context.Context, productId uint64) error {
	// First get all SKUs to invalidate cache
	skus, err := m.FindManyByProductId(ctx, productId)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(skus)*2)
	for _, sku := range skus {
		keys = append(keys,
			fmt.Sprintf("%s%v", cacheMallProductSkusIdPrefix, sku.Id),
			fmt.Sprintf("%s%v", cacheMallProductSkusSkuCodePrefix, sku.SkuCode))
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `product_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, productId)
	}, keys...)

	return err
}
