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

func TestUpdateItemLogic_UpdateItem(t *testing.T) {
    configFile := flag.String("f", "../../etc/cart.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    mockProduct := NewProductService(t)
    ctx.ProductRpc = mockProduct

    // Create test data
    testItem := &model.CartItems{
        UserId:      1,
        ProductId:   1,
        SkuId:      1,
        ProductName: "Test Product",
        SkuName:    "SKU1",
        Price:      100,
        Quantity:   1,
        Selected:   1,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    result, err := ctx.CartItemsModel.Insert(context.Background(), testItem)
    assert.NoError(t, err)

    testStats := &model.CartStatistics{
        UserId:           1,
        TotalQuantity:    1,
        SelectedQuantity: 1,
        TotalAmount:      100,
        SelectedAmount:   100,
    }
    err = ctx.CartStatsModel.Upsert(context.Background(), testStats)
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *cart.UpdateItemRequest
        mock    func()
        wantErr error
        check   func(*testing.T)
    }{
        {
            name: "Update quantity successfully",
            req: &cart.UpdateItemRequest{
                UserId:    1,
                ProductId: 1,
                SkuId:    1,
                Quantity: 2,
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
            check: func(t *testing.T) {
                item, err := ctx.CartItemsModel.FindOneByUserIdSkuId(context.Background(), 1, 1)
                assert.NoError(t, err)
                assert.Equal(t, int64(2), item.Quantity)

                stats, err := ctx.CartStatsModel.FindOne(context.Background(), 1)
                assert.NoError(t, err)
                assert.Equal(t, int64(2), stats.TotalQuantity)
                assert.Equal(t, float64(200), stats.TotalAmount)
            },
        },
        {
            name: "Invalid quantity",
            req: &cart.UpdateItemRequest{
                UserId:    1,
                ProductId: 1,
                SkuId:    1,
                Quantity: 0,
            },
            mock:    func() {},
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Out of stock",
            req: &cart.UpdateItemRequest{
                UserId:    1,
                ProductId: 1,
                SkuId:    1,
                Quantity: 20,
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
            wantErr: zeroerr.ErrItemOutOfStock,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewUpdateItemLogic(context.Background(), ctx)
            resp, err := l.UpdateItem(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.True(t, resp.Success)
                if tt.check != nil {
                    tt.check(t)
                }
            }
        })
    }

    // Cleanup
    id, _ := result.LastInsertId()
    _ = ctx.CartItemsModel.Delete(context.Background(), uint64(id))
    _ = ctx.CartStatsModel.Delete(context.Background(), 1)
}