package logic

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestCheckStockLogic_CheckStock(t *testing.T) {
    configFile := flag.String("f", "../../etc/cart.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    mockProduct := NewProductService(t)
    ctx.ProductRpc = mockProduct

    // Create test cart items
    testItems := []*model.CartItems{
        {
            UserId:      1,
            ProductId:   1,
            SkuId:      1,
            ProductName: "Test Product 1",
            SkuName:    "SKU1",
            Price:      100,
            Quantity:   2,
            Selected:   1,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        },
    }

    err := ctx.CartItemsModel.BatchInsert(context.Background(), testItems)
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *cart.CheckStockRequest
        mock    func()
        wantErr error
        check   func(*testing.T, *cart.CheckStockResponse)
    }{
        {
            name: "All items in stock",
            req: &cart.CheckStockRequest{
                UserId: 1,
            },
            mock: func() {
                mockProduct.EXPECT().GetSku(context.Background(), &product.GetSkuRequest{
                    Id: 1,
                }).Return(&product.GetSkuResponse{
                    Sku: &product.Sku{
                        Id:    1,
                        Stock: 10,
                    },
                }, nil)
            },
            wantErr: nil,
            check: func(t *testing.T, resp *cart.CheckStockResponse) {
                assert.True(t, resp.AllInStock)
                assert.Empty(t, resp.OutOfStockItems)
            },
        },
        {
            name: "Invalid user ID",
            req: &cart.CheckStockRequest{
                UserId: 0,
            },
            mock:    func() {},
            wantErr: zeroerr.ErrInvalidParam,
            check:   nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewCheckStockLogic(context.Background(), ctx)
            resp, err := l.CheckStock(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                if tt.check != nil {
                    tt.check(t, resp)
                }
            }
        })
    }

    // Cleanup
    err = ctx.CartItemsModel.DeleteByUserId(context.Background(), 1)
    assert.NoError(t, err)
}