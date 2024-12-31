package logic

import (
	"context"
	"encoding/json"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestUpdateSkuStockLogic_UpdateSkuStock(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test SKU
	attrs := []map[string]string{
		{"color": "red", "size": "L"},
	}
	attrsJSON, err := json.Marshal(attrs)
	assert.NoError(t, err)

	testSku := &model.Skus{
		ProductId: 1,
		SkuCode:   "TEST-SKU-001",
		Attributes: string(attrsJSON),
		Price: 9999,
		Stock: 100,
		Sales: 0,
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
		req     *product.UpdateSkuStockRequest
		wantErr error
	}{
		{
			name: "Increase stock",
			req: &product.UpdateSkuStockRequest{
				Id:        skuId,
				Increment: 50,
			},
			wantErr: nil,
		},
		{
			name: "Decrease stock",
			req: &product.UpdateSkuStockRequest{
				Id:        skuId,
				Increment: -30,
			},
			wantErr: nil,
		},
		{
			name: "Invalid SKU ID",
			req: &product.UpdateSkuStockRequest{
				Id:        0,
				Increment: 10,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent SKU",
			req: &product.UpdateSkuStockRequest{
				Id:        99999,
				Increment: 10,
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
		{
			name: "Invalid negative stock",
			req: &product.UpdateSkuStockRequest{
				Id:        skuId,
				Increment: -200,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateSkuStockLogic(context.Background(), ctx)
			resp, err := l.UpdateSkuStock(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify stock update
				updated, err := ctx.SkusModel.FindOne(context.Background(), uint64(skuId))
				assert.NoError(t, err)
				assert.Equal(t, testSku.Stock+tt.req.Increment, updated.Stock)
				testSku.Stock = updated.Stock // Update for next test case
			}
		})
	}
}
