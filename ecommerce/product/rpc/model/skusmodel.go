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
		BatchInsert(ctx context.Context, data []*Skus) error
		DeleteByProductId(ctx context.Context, productId uint64) error
		UpdateSkus(ctx context.Context, id uint64, updates map[string]interface{}) error
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

func (m *customSkusModel) BatchInsert(ctx context.Context, data []*Skus) error {
	if len(data) == 0 {
		return nil
	}

	values := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data)*6)
	for _, sku := range data {
		values = append(values, "(?, ?, ?, ?, ?)")
		args = append(args, sku.ProductId, sku.SkuCode, sku.Attributes, sku.Price, sku.Sales)
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

func (m *customSkusModel) UpdateSkus(ctx context.Context, id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	var sets []string
	var args []interface{}
	for k, v := range updates {
		// Handle increments
		if k == "stock" || k == "sales" {
			sets = append(sets, fmt.Sprintf("`%s` = `%s` + ?", k, k))
		} else {
			sets = append(sets, fmt.Sprintf("`%s` = ?", k))
		}
		args = append(args, v)
	}
	args = append(args, id)

	mallProductSkusIdKey := fmt.Sprintf("%s%v", cacheMallProductSkusIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(sets, ", "))
		return conn.ExecCtx(ctx, query, args...)
	}, mallProductSkusIdKey)

	return err
}