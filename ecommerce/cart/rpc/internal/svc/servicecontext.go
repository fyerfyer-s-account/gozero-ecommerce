package svc

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config            config.Config
    CartItemsModel    model.CartItemsModel
    CartStatsModel    model.CartStatisticsModel
    ProductRpc        productservice.ProductService
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    
    return &ServiceContext{
        Config:            c,
        CartItemsModel:    model.NewCartItemsModel(conn, c.CacheRedis),
        CartStatsModel:    model.NewCartStatisticsModel(conn, c.CacheRedis),
        ProductRpc:        productservice.NewProductService(zrpc.MustNewClient(c.ProductRpc)),
    }
}