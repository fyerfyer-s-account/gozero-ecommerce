package logic

import (
	"context"
	"database/sql"
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

func TestUnselectAllLogic_UnselectAll(t *testing.T) {
    configFile := flag.String("f", "../../etc/cart.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test data
    testItems := []*model.CartItems{
        {
            UserId:      1,
            ProductId:   1,
            SkuId:      1,
            ProductName: "Product 1",
            SkuName:    "SKU1",
            Image: sql.NullString{
                String: "test1.jpg",
                Valid:  true,
            },
            Price:     100,
            Quantity:  2,
            Selected:  1,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        {
            UserId:      1,
            ProductId:   2,
            SkuId:      2,
            ProductName: "Product 2",
            SkuName:    "SKU2",
            Price:      200,
            Quantity:   1,
            Selected:   1,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        },
    }

    err := ctx.CartItemsModel.BatchInsert(context.Background(), testItems)
    assert.NoError(t, err)

    testStats := &model.CartStatistics{
        UserId:           1,
        TotalQuantity:    3,
        SelectedQuantity: 3,
        TotalAmount:      400,
        SelectedAmount:   400,
    }
    err = ctx.CartStatsModel.Upsert(context.Background(), testStats)
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *cart.UnselectAllRequest
        wantErr error
        check   func(*testing.T)
    }{
        {
            name: "Unselect all successfully",
            req: &cart.UnselectAllRequest{
                UserId: 1,
            },
            wantErr: nil,
            check: func(t *testing.T) {
                items, err := ctx.CartItemsModel.FindByUserId(context.Background(), 1)
                assert.NoError(t, err)
                for _, item := range items {
                    assert.Equal(t, int64(0), item.Selected)
                }

                stats, err := ctx.CartStatsModel.FindOne(context.Background(), 1)
                assert.NoError(t, err)
                assert.Equal(t, int64(0), stats.SelectedQuantity)
                assert.Equal(t, float64(0), stats.SelectedAmount)
            },
        },
        {
            name: "Invalid user ID",
            req: &cart.UnselectAllRequest{
                UserId: 0,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Empty cart",
            req: &cart.UnselectAllRequest{
                UserId: 999,
            },
            wantErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewUnselectAllLogic(context.Background(), ctx)
            resp, err := l.UnselectAll(tt.req)

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
    err = ctx.CartItemsModel.DeleteByUserId(context.Background(), 1)
    assert.NoError(t, err)
    err = ctx.CartStatsModel.Delete(context.Background(), 1)
    assert.NoError(t, err)
}