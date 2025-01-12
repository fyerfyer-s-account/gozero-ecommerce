package svc

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config         config.Config
    StocksModel    model.StocksModel
    StockLocksModel    model.StockLocksModel
    StockRecordsModel  model.StockRecordsModel
    WarehousesModel    model.WarehousesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    
    return &ServiceContext{
        Config:         c,
        StocksModel:    model.NewStocksModel(conn, c.CacheRedis),
        StockLocksModel:    model.NewStockLocksModel(conn, c.CacheRedis),
        StockRecordsModel:  model.NewStockRecordsModel(conn, c.CacheRedis),
        WarehousesModel:    model.NewWarehousesModel(conn, c.CacheRedis),
    }
}