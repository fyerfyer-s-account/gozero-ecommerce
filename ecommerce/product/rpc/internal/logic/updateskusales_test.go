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

func TestUpdateSkuSalesLogic_UpdateSkuSales(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create attributes JSON
	attrs := []map[string]string{
		{"color": "red", "size": "L"},
	}
	attrsJSON, err := json.Marshal(attrs)
	assert.NoError(t, err)

	// Create test SKU
	testSku := &model.Skus{
		ProductId:  1,
		SkuCode:    "TEST-SKU-001",
		Attributes: string(attrsJSON),
		Price:      9999,
		Stock:      100,
		Sales:      50,
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
		req     *product.UpdateSkuSalesRequest
		wantErr error
	}{
		{
			name: "Increment sales",
			req: &product.UpdateSkuSalesRequest{
				Id:        skuId,
				Increment: 10,
			},
			wantErr: nil,
		},
		{
			name: "Decrement sales",
			req: &product.UpdateSkuSalesRequest{
				Id:        skuId,
				Increment: -5,
			},
			wantErr: nil,
		},
		{
			name: "Invalid SKU ID",
			req: &product.UpdateSkuSalesRequest{
				Id:        0,
				Increment: 10,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent SKU",
			req: &product.UpdateSkuSalesRequest{
				Id:        99999,
				Increment: 10,
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateSkuSalesLogic(context.Background(), ctx)
			resp, err := l.UpdateSkuSales(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify sales update
				updated, err := ctx.SkusModel.FindOne(context.Background(), uint64(skuId))
				assert.NoError(t, err)
				if tt.req.Increment > 0 {
					assert.Equal(t, testSku.Sales+tt.req.Increment, updated.Sales)
				} else {
					assert.Equal(t, testSku.Sales-(-tt.req.Increment), updated.Sales)
				}
				testSku.Sales = updated.Sales // Update for next test case
			}
		})
	}
}
