package svc

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	ProductsModel       model.ProductsModel
	CategoriesModel     model.CategoriesModel
	SkusModel           model.SkusModel
	ProductReviewsModel model.ProductReviewsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:              c,
		ProductsModel:       model.NewProductsModel(conn, c.CacheRedis),
		CategoriesModel:     model.NewCategoriesModel(conn, c.CacheRedis),
		SkusModel:           model.NewSkusModel(conn, c.CacheRedis),
		ProductReviewsModel: model.NewProductReviewsModel(conn, c.CacheRedis),
	}
}
