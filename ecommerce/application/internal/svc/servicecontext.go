package svc

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/middleware"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cartclient"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketingclient"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/orderservice"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/paymentclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"
	// "github.com/fyerfyer/gozero-ecommerce/ecommerce/search/rpc/searchclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Auth      rest.Middleware
	AdminAuth rest.Middleware

	Redis *redis.Redis

	UserRpc    userclient.User
	ProductRpc productservice.ProductService
	// CartRpc      cartclient.Cart
	// OrderRpc     orderservice.OrderService
	// PaymentRpc   paymentclient.Payment
	// InventoryRpc inventoryclient.Inventory
	// MarketingRpc marketingclient.Marketing
	// SearchRpc    searchclient.Search
	// MessageRpc   messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize Redis
	rdb := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = c.Redis.Type
		r.Pass = c.Redis.Pass
	})

	return &ServiceContext{
		Config: c,
		Redis:  rdb,

		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc: productservice.NewProductService(zrpc.MustNewClient(c.ProductRpc)),
		AdminAuth:  middleware.NewAdminAuthMiddleware(c).Handle,
		Auth:       middleware.NewAuthMiddleware(c).Handle,
		// CartRpc:      cartclient.NewCart(zrpc.MustNewClient(c.CartRpc)),
		// OrderRpc:     orderservice.NewOrderService(zrpc.MustNewClient(c.OrderRpc)),
		// PaymentRpc:   paymentclient.NewPayment(zrpc.MustNewClient(c.PaymentRpc)),
		// InventoryRpc: inventoryclient.NewInventory(zrpc.MustNewClient(c.InventoryRpc)),
		// MarketingRpc: marketingclient.NewMarketing(zrpc.MustNewClient(c.MarketingRpc)),
		// SearchRpc:    searchclient.NewSearch(zrpc.MustNewClient(c.SearchRpc)),
		// MessageRpc:   messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc)),
	}
}
