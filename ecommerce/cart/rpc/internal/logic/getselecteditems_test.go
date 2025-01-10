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

func TestGetSelectedItemsLogic_GetSelectedItems(t *testing.T) {
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
            ProductName: "Selected Product",
            SkuName:    "SKU1",
            Image: sql.NullString{
                String: "test.jpg",
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
            ProductName: "Unselected Product",
            SkuName:    "SKU2",
            Price:      200,
            Quantity:   1,
            Selected:   0,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        },
    }

    err := ctx.CartItemsModel.BatchInsert(context.Background(), testItems)
    assert.NoError(t, err)

    testStats := &model.CartStatistics{
        UserId:           1,
        TotalQuantity:    3,
        SelectedQuantity: 2,
        TotalAmount:      300,
        SelectedAmount:   200,
    }
    err = ctx.CartStatsModel.Upsert(context.Background(), testStats)
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *cart.GetSelectedItemsRequest
        wantErr error
        check   func(*testing.T, *cart.GetSelectedItemsResponse)
    }{
        {
            name: "Get selected items successfully",
            req: &cart.GetSelectedItemsRequest{
                UserId: 1,
            },
            wantErr: nil,
            check: func(t *testing.T, resp *cart.GetSelectedItemsResponse) {
                assert.Len(t, resp.Items, 1)
                assert.Equal(t, "Selected Product", resp.Items[0].ProductName)
                assert.Equal(t, int32(2), resp.TotalQuantity)
                assert.Equal(t, float64(200), resp.TotalPrice)
            },
        },
        {
            name: "Empty cart",
            req: &cart.GetSelectedItemsRequest{
                UserId: 999,
            },
            wantErr: nil,
            check: func(t *testing.T, resp *cart.GetSelectedItemsResponse) {
                assert.Empty(t, resp.Items)
                assert.Equal(t, int32(0), resp.TotalQuantity)
                assert.Equal(t, float64(0), resp.TotalPrice)
            },
        },
        {
            name: "Invalid user ID",
            req: &cart.GetSelectedItemsRequest{
                UserId: 0,
            },
            wantErr: zeroerr.ErrInvalidParam,
            check:   nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewGetSelectedItemsLogic(context.Background(), ctx)
            resp, err := l.GetSelectedItems(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                if tt.check != nil {
                    tt.check(t, resp)
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