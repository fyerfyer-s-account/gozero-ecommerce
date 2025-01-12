package model

import (
	"context"
    "database/sql"
    "fmt"
	"strings"
	
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StockRecordsModel = (*customStockRecordsModel)(nil)

type (
	// StockRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockRecordsModel.
	StockRecordsModel interface {
		stockRecordsModel
		// Transaction support
        Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
        // Batch operations
        BatchInsert(ctx context.Context, records []*StockRecords) error
        // Query methods
        FindByOrderNo(ctx context.Context, orderNo string) ([]*StockRecords, error)
        FindBySkuAndWarehouse(ctx context.Context, skuId, warehouseId uint64) ([]*StockRecords, error)
        // List with pagination
        List(ctx context.Context, skuId, warehouseId uint64, recordType int32, page, pageSize int32) ([]*StockRecords, int64, error)
	}

	customStockRecordsModel struct {
		*defaultStockRecordsModel
	}
)

// NewStockRecordsModel returns a model for the database table.
func NewStockRecordsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StockRecordsModel {
	return &customStockRecordsModel{
		defaultStockRecordsModel: newStockRecordsModel(conn, c, opts...),
	}
}

func (m *customStockRecordsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, fn)
}

func (m *customStockRecordsModel) BatchInsert(ctx context.Context, records []*StockRecords) error {
    if len(records) == 0 {
        return nil
    }

    values := make([]string, 0, len(records))
    args := make([]interface{}, 0, len(records)*7)
    
    for _, record := range records {
        values = append(values, "(?, ?, ?, ?, ?, ?, ?)")
        args = append(args, 
            record.SkuId,
            record.WarehouseId,
            record.Type,
            record.Quantity,
            record.OrderNo,
            record.Remark,
            record.Operator,
        )
    }

    query := fmt.Sprintf("insert into %s (%s) values %s",
        m.table, stockRecordsRowsExpectAutoSet, strings.Join(values, ","))
    
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })
    return err
}

func (m *customStockRecordsModel) FindByOrderNo(ctx context.Context, orderNo string) ([]*StockRecords, error) {
    var records []*StockRecords
    query := fmt.Sprintf("select %s from %s where `order_no` = ?", stockRecordsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, orderNo)
    return records, err
}

func (m *customStockRecordsModel) FindBySkuAndWarehouse(ctx context.Context, skuId, warehouseId uint64) ([]*StockRecords, error) {
    var records []*StockRecords
    query := fmt.Sprintf("select %s from %s where `sku_id` = ? and `warehouse_id` = ?", 
        stockRecordsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &records, query, skuId, warehouseId)
    return records, err
}

func (m *customStockRecordsModel) List(ctx context.Context, skuId, warehouseId uint64, recordType int32, page, pageSize int32) ([]*StockRecords, int64, error) {
    // Build where clause
    conditions := []string{}
    args := []interface{}{}
    
    if skuId > 0 {
        conditions = append(conditions, "`sku_id` = ?")
        args = append(args, skuId)
    }
    if warehouseId > 0 {
        conditions = append(conditions, "`warehouse_id` = ?")
        args = append(args, warehouseId)
    }
    if recordType > 0 {
        conditions = append(conditions, "`type` = ?")
        args = append(args, recordType)
    }
    
    whereClause := ""
    if len(conditions) > 0 {
        whereClause = "where " + strings.Join(conditions, " and ")
    }

    // Query records
    query := fmt.Sprintf("select %s from %s %s order by id desc limit ?, ?", 
        stockRecordsRows, m.table, whereClause)
    countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, whereClause)
    
    var count int64
    err := m.QueryRowNoCacheCtx(ctx, &count, countQuery, args...)
    if err != nil {
        return nil, 0, err
    }

    offset := (page - 1) * pageSize
    args = append(args, offset, pageSize)
    
    var records []*StockRecords
    err = m.QueryRowsNoCacheCtx(ctx, &records, query, args...)
    if err != nil {
        return nil, 0, err
    }

    return records, count, nil
}