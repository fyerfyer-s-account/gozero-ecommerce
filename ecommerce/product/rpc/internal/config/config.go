package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	// Product specific settings
	MaxCategoryLevel    int `json:",default=3"`
	MaxSkusPerProduct   int `json:",default=100"`
	MaxImagesPerProduct int `json:",default=10"`
	MaxReviewImages     int `json:",default=5"`
	MinReviewLength     int `json:",default=5"`
	MaxReviewLength     int `json:",default=500"`
	PageSize            int `json:",default=10"`
}
