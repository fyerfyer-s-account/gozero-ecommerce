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

func TestUpdateSkuPriceLogic_UpdateSkuPrice(t *testing.T) {
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
		Sales:      0,
	}

	var skuId int64
	result, err := ctx.SkusModel.Insert(context.Background(), testSku)
	if err != nil {
		t.Fatalf("Failed to create test SKU: %v", err)
	} else {
		skuId, err = result.LastInsertId()
		assert.NoError(t, err)
	}

	// Cleanup after test
	defer func() {
		if skuId > 0 {
			_ = ctx.SkusModel.Delete(context.Background(), uint64(skuId))
		}
	}()

	tests := []struct {
		name    string
		req     *product.UpdateSkuPriceRequest
		wantErr error
	}{
		{
			name: "Valid price update",
			req: &product.UpdateSkuPriceRequest{
				Id:    skuId,
				Price: 19999, // 199.99
			},
			wantErr: nil,
		},
		{
			name: "Invalid SKU ID",
			req: &product.UpdateSkuPriceRequest{
				Id:    0,
				Price: 9999,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Negative price",
			req: &product.UpdateSkuPriceRequest{
				Id:    skuId,
				Price: -1000,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent SKU",
			req: &product.UpdateSkuPriceRequest{
				Id:    99999,
				Price: 9999,
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateSkuPriceLogic(context.Background(), ctx)
			resp, err := l.UpdateSkuPrice(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify price update
				updated, err := ctx.SkusModel.FindOne(context.Background(), uint64(skuId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Price, updated.Price)
			}
		})
	}
}
