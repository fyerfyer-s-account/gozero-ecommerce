// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	stocksFieldNames          = builder.RawFieldNames(&Stocks{})
	stocksRows                = strings.Join(stocksFieldNames, ",")
	stocksRowsExpectAutoSet   = strings.Join(stringx.Remove(stocksFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	stocksRowsWithPlaceHolder = strings.Join(stringx.Remove(stocksFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
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
		conn  sqlx.SqlConn
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

func newStocksModel(conn sqlx.SqlConn) *defaultStocksModel {
	return &defaultStocksModel{
		conn:  conn,
		table: "`stocks`",
	}
}

func (m *defaultStocksModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultStocksModel) FindOne(ctx context.Context, id uint64) (*Stocks, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", stocksRows, m.table)
	var resp Stocks
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStocksModel) FindOneBySkuIdWarehouseId(ctx context.Context, skuId uint64, warehouseId uint64) (*Stocks, error) {
	var resp Stocks
	query := fmt.Sprintf("select %s from %s where `sku_id` = ? and `warehouse_id` = ? limit 1", stocksRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, skuId, warehouseId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStocksModel) Insert(ctx context.Context, data *Stocks) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, stocksRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.SkuId, data.WarehouseId, data.Available, data.Locked, data.Total, data.AlertQuantity)
	return ret, err
}

func (m *defaultStocksModel) Update(ctx context.Context, newData *Stocks) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, stocksRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, newData.SkuId, newData.WarehouseId, newData.Available, newData.Locked, newData.Total, newData.AlertQuantity, newData.Id)
	return err
}

func (m *defaultStocksModel) tableName() string {
	return m.table
}
