package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestUpdateSkuLogic_UpdateSku(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test SKU
	testSku := &model.Skus{
		ProductId: 1,
		SkuCode:   "TEST001",
		Price:     99.99,
		Stock:     100,
		Sales:     10,
		Attributes: sql.NullString{
			String: `[{"key":"color","value":"red"}]`,
			Valid:  true,
		},
	}

	result, err := ctx.SkusModel.Insert(context.Background(), testSku)
	assert.NoError(t, err)
	skuId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.SkusModel.Delete(context.Background(), uint64(skuId))
	}()

	tests := []struct {
		name    string
		req     *product.UpdateSkuRequest
		wantErr bool
		verify  func(*testing.T, uint64, error)
	}{
		{
			name: "Update price only",
			req: &product.UpdateSkuRequest{
				Id:    skuId,
				Price: 199.99,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				time.Sleep(100 * time.Millisecond) 
				updated, err := ctx.SkusModel.FindOne(context.Background(), id)
				assert.NoError(t, err)
				assert.Equal(t, 199.99, updated.Price)
			},
		},
		{
			name: "Update stock increment",
			req: &product.UpdateSkuRequest{
				Id:             skuId,
				StockIncrement: 50,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				time.Sleep(100 * time.Millisecond)
				updated, err := ctx.SkusModel.FindOne(context.Background(), id)
				assert.NoError(t, err)
				assert.Equal(t, int64(150), updated.Stock) // 100 + 50
			},
		},
		{
			name: "Update sales increment",
			req: &product.UpdateSkuRequest{
				Id:             skuId,
				SalesIncrement: 5,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				time.Sleep(100 * time.Millisecond) 
				updated, err := ctx.SkusModel.FindOne(context.Background(), id)
				assert.NoError(t, err)
				assert.Equal(t, int64(15), updated.Sales) // 10 + 5
			},
		},
		{
			name: "Invalid SKU ID",
			req: &product.UpdateSkuRequest{
				Id:    0,
				Price: 299.99,
			},
			wantErr: true,
			verify: func(t *testing.T, id uint64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateSkuLogic(context.Background(), ctx)
			resp, err := l.UpdateSku(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
				if tt.verify != nil {
					tt.verify(t, uint64(tt.req.Id), err)
				}
			}
		})
	}
}
