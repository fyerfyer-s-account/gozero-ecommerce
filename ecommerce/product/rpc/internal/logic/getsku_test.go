package logic

import (
	"context"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestGetSkuLogic_GetSku(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/product.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Prepare attributes JSON
	attrs := []*product.SkuAttribute{
		{Key: "color", Value: "red"},
		{Key: "size", Value: "M"},
	}
	attrsJSON, err := json.Marshal(attrs)
	assert.NoError(t, err)

	// Create test SKU
	testSku := &model.Skus{
		ProductId:  1,
		SkuCode:    "TEST-SKU-001",
		Price:      99.99,
		Stock:      100,
		Attributes: string(attrsJSON),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err := ctx.SkusModel.Insert(context.Background(), testSku)
	assert.NoError(t, err)
	skuId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		err := ctx.SkusModel.Delete(context.Background(), uint64(skuId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *product.GetSkuRequest
		wantErr error
	}{
		{
			name: "Valid SKU",
			req: &product.GetSkuRequest{
				Id: skuId,
			},
			wantErr: nil,
		},
		{
			name: "Invalid SKU ID",
			req: &product.GetSkuRequest{
				Id: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent SKU",
			req: &product.GetSkuRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetSkuLogic(context.Background(), ctx)
			resp, err := l.GetSku(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Sku)
				assert.Equal(t, testSku.SkuCode, resp.Sku.SkuCode)
				assert.Equal(t, testSku.Price, resp.Sku.Price)
				assert.Equal(t, int64(testSku.Stock), resp.Sku.Stock)
				assert.Len(t, resp.Sku.Attributes, 2)
			}
		})
	}
}
