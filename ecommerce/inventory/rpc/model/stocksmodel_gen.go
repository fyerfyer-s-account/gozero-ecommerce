// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	stocksFieldNames          = builder.RawFieldNames(&Stocks{})
	stocksRows                = strings.Join(stocksFieldNames, ",")
	stocksRowsExpectAutoSet   = strings.Join(stringx.Remove(stocksFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	stocksRowsWithPlaceHolder = strings.Join(stringx.Remove(stocksFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallInventoryStocksIdPrefix               = "cache:mallInventory:stocks:id:"
	cacheMallInventoryStocksSkuIdWarehouseIdPrefix = "cache:mallInventory:stocks:skuId:warehouseId:"
)

type (
	stocksModel interface {
		Insert(ctx context.Context, data *Stocks) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Stocks, error)
		FindOneBySkuIdWarehouseId(ctx context.Context, skuId uint64, warehouseId uint64) (*Stocks, error)
		Update(ctx context.Context, data *Stocks) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultStocksModel struct {
		sqlc.CachedConn
		table string
	}

	Stocks struct {
		Id            uint64    `db:"id"`             // åº“å­˜ID
		SkuId         uint64    `db:"sku_id"`         // SKU ID
		WarehouseId   uint64    `db:"warehouse_id"`   // ä»“åº“ID
		Available     int64     `db:"available"`      // å¯ç”¨åº“å­˜
		Locked        int64     `db:"locked"`         // é”å®šåº“å­˜
		Total         int64     `db:"total"`          // æ€»åº“å­˜
		AlertQuantity int64     `db:"alert_quantity"` // åº“å­˜é¢„è­¦æ•°é‡
		CreatedAt     time.Time `db:"created_at"`     // åˆ›å»ºæ—¶é—´
		UpdatedAt     time.Time `db:"updated_at"`     // æ›´æ–°æ—¶é—´
	}
)

func newStocksModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultStocksModel {
	return &defaultStocksModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`stocks`",
	}
}

func (m *defaultStocksModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallInventoryStocksIdKey := fmt.Sprintf("%s%v", cacheMallInventoryStocksIdPrefix, id)
	mallInventoryStocksSkuIdWarehouseIdKey := fmt.Sprintf("%s%v:%v", cacheMallInventoryStocksSkuIdWarehouseIdPrefix, data.SkuId, data.WarehouseId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mallInventoryStocksIdKey, mallInventoryStocksSkuIdWarehouseIdKey)
	return err
}

func (m *defaultStocksModel) FindOne(ctx context.Context, id uint64) (*Stocks, error) {
	mallInventoryStocksIdKey := fmt.Sprintf("%s%v", cacheMallInventoryStocksIdPrefix, id)
	var resp Stocks
	err := m.QueryRowCtx(ctx, &resp, mallInventoryStocksIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", stocksRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStocksModel) FindOneBySkuIdWarehouseId(ctx context.Context, skuId uint64, warehouseId uint64) (*Stocks, error) {
	mallInventoryStocksSkuIdWarehouseIdKey := fmt.Sprintf("%s%v:%v", cacheMallInventoryStocksSkuIdWarehouseIdPrefix, skuId, warehouseId)
	var resp Stocks
	err := m.QueryRowIndexCtx(ctx, &resp, mallInventoryStocksSkuIdWarehouseIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `sku_id` = ? and `warehouse_id` = ? limit 1", stocksRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, skuId, warehouseId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStocksModel) Insert(ctx context.Context, data *Stocks) (sql.Result, error) {
	mallInventoryStocksIdKey := fmt.Sprintf("%s%v", cacheMallInventoryStocksIdPrefix, data.Id)
	mallInventoryStocksSkuIdWarehouseIdKey := fmt.Sprintf("%s%v:%v", cacheMallInventoryStocksSkuIdWarehouseIdPrefix, data.SkuId, data.WarehouseId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, stocksRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.SkuId, data.WarehouseId, data.Available, data.Locked, data.Total, data.AlertQuantity)
	}, mallInventoryStocksIdKey, mallInventoryStocksSkuIdWarehouseIdKey)
	return ret, err
}

func (m *defaultStocksModel) Update(ctx context.Context, newData *Stocks) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	mallInventoryStocksIdKey := fmt.Sprintf("%s%v", cacheMallInventoryStocksIdPrefix, data.Id)
	mallInventoryStocksSkuIdWarehouseIdKey := fmt.Sprintf("%s%v:%v", cacheMallInventoryStocksSkuIdWarehouseIdPrefix, data.SkuId, data.WarehouseId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, stocksRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.SkuId, newData.WarehouseId, newData.Available, newData.Locked, newData.Total, newData.AlertQuantity, newData.Id)
	}, mallInventoryStocksIdKey, mallInventoryStocksSkuIdWarehouseIdKey)
	return err
}

func (m *defaultStocksModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallInventoryStocksIdPrefix, primary)
}

func (m *defaultStocksModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", stocksRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultStocksModel) tableName() string {
	return m.table
}
