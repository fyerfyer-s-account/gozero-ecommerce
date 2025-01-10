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

func TestSelectItemLogic_SelectItem(t *testing.T) {
    configFile := flag.String("f", "../../etc/cart.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test data
    testItem := &model.CartItems{
        UserId:      1,
        ProductId:   1,
        SkuId:      1,
        ProductName: "Test Product",
        SkuName:    "SKU1",
        Image: sql.NullString{
            String: "test.jpg",
            Valid:  true,
        },
        Price:     100,
        Quantity:  2,
        Selected:  0,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    result, err := ctx.CartItemsModel.Insert(context.Background(), testItem)
    assert.NoError(t, err)

    testStats := &model.CartStatistics{
        UserId:           1,
        TotalQuantity:    2,
        SelectedQuantity: 0,
        TotalAmount:      200,
        SelectedAmount:   0,
    }
    err = ctx.CartStatsModel.Upsert(context.Background(), testStats)
    assert.NoError(t, err)

    tests := []struct {
        name    string
        req     *cart.SelectItemRequest
        wantErr error
        check   func(*testing.T)
    }{
        {
            name: "Select item successfully",
            req: &cart.SelectItemRequest{
                UserId:    1,
                ProductId: 1,
                SkuId:    1,
            },
            wantErr: nil,
            check: func(t *testing.T) {
                item, err := ctx.CartItemsModel.FindOneByUserIdSkuId(context.Background(), 1, 1)
                assert.NoError(t, err)
                assert.Equal(t, int64(1), item.Selected)

                stats, err := ctx.CartStatsModel.FindOne(context.Background(), 1)
                assert.NoError(t, err)
                assert.Equal(t, int64(2), stats.SelectedQuantity)
                assert.Equal(t, float64(200), stats.SelectedAmount)
            },
        },
        {
            name: "Invalid user ID",
            req: &cart.SelectItemRequest{
                UserId:    0,
                ProductId: 1,
                SkuId:    1,
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "Item not found",
            req: &cart.SelectItemRequest{
                UserId:    999,
                ProductId: 999,
                SkuId:    999,
            },
            wantErr: zeroerr.ErrSelectFailed,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewSelectItemLogic(context.Background(), ctx)
            resp, err := l.SelectItem(tt.req)

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