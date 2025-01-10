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
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestClearCartLogic_ClearCart(t *testing.T) {
    configFile := flag.String("f", "../../etc/cart.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test cart items
    testItems := []*model.CartItems{
        {
            UserId:      1,
            ProductId:   1,
            SkuId:      1,
            ProductName: "Test Product",
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

    // Create test cart statistics
    testStats := &model.CartStatistics{
        UserId:           1,
        TotalQuantity:    2,
        SelectedQuantity: 2,
        TotalAmount:      200,
        SelectedAmount:   200,
    }
    err = ctx.CartStatsModel.Upsert(context.Background(), testStats)
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *cart.ClearCartRequest
        wantErr error
    }{
        {
            name: "Successfully clear cart",
            req: &cart.ClearCartRequest{
                UserId: 1,
            },
            wantErr: nil,
        },
        {
            name: "Invalid user ID",
            req: &cart.ClearCartRequest{
                UserId: 0,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Empty cart",
            req: &cart.ClearCartRequest{
                UserId: 999,
            },
            wantErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewClearCartLogic(context.Background(), ctx)
            resp, err := l.ClearCart(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.True(t, resp.Success)

                // Verify cart is empty
                items, err := ctx.CartItemsModel.FindByUserId(context.Background(), uint64(tt.req.UserId))
                assert.NoError(t, err)
                assert.Empty(t, items)

                // Verify statistics are reset
                stats, err := ctx.CartStatsModel.FindOne(context.Background(), uint64(tt.req.UserId))
                if err != model.ErrNotFound {
                    assert.NoError(t, err)
                    assert.Equal(t, int64(0), stats.TotalQuantity)
                    assert.Equal(t, int64(0), stats.SelectedQuantity)
                    assert.Equal(t, float64(0), stats.TotalAmount)
                    assert.Equal(t, float64(0), stats.SelectedAmount)
                }
            }
        })
    }
}