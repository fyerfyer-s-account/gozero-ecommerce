package svc

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	// Models
	UsersModel              model.UsersModel
	UserAddressesModel      model.UserAddressesModel
	LoginRecordsModel       model.LoginRecordsModel
	WalletAccountsModel     model.WalletAccountsModel
	WalletTransactionsModel model.WalletTransactionsModel
	AdminModel              model.AdminsModel

	// Redis clients
	BizRedis     *redis.Redis
	RefreshRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	bizRedis := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
		r.Type = c.CacheRedis[0].Type
		r.Pass = c.CacheRedis[0].Pass
	})

	refreshRedis := redis.New(c.JwtAuth.RefreshRedis.Host, func(r *redis.Redis) {
		r.Type = c.JwtAuth.RefreshRedis.Type
		r.Pass = c.JwtAuth.RefreshRedis.Pass
	})

	return &ServiceContext{
		Config: c,

		AdminModel:              model.NewAdminsModel(sqlConn),
		UsersModel:              model.NewUsersModel(sqlConn),
		UserAddressesModel:      model.NewUserAddressesModel(sqlConn),
		LoginRecordsModel:       model.NewLoginRecordsModel(sqlConn),
		WalletAccountsModel:     model.NewWalletAccountsModel(sqlConn),
		WalletTransactionsModel: model.NewWalletTransactionsModel(sqlConn),

		BizRedis:     bizRedis,
		RefreshRedis: refreshRedis,
	}
}
