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

func TestDeleteSkuLogic_DeleteSku(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product
	testProduct := &model.Products{
		Name:   "Test Product",
		Price:  99.99,
		Status: 1,
	}
	result, err := ctx.ProductsModel.Insert(context.Background(), testProduct)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test SKU with valid attributes
	attrs := []map[string]string{
		{"key": "color", "value": "red"},
		{"key": "size", "value": "M"},
	}
	attrsJSON, err := json.Marshal(attrs)
	assert.NoError(t, err)

	testSku := &model.Skus{
		ProductId:  uint64(productId),
		SkuCode:    "TEST-SKU-001",
		Price:      99.99,
		Attributes: string(attrsJSON),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err = ctx.SkusModel.Insert(context.Background(), testSku)
	assert.NoError(t, err)
	skuId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name    string
		req     *product.DeleteSkuRequest
		wantErr error
	}{
		{
			name: "Valid SKU",
			req: &product.DeleteSkuRequest{
				Id: skuId,
			},
			wantErr: nil,
		},
		{
			name: "Invalid SKU ID",
			req: &product.DeleteSkuRequest{
				Id: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent SKU",
			req: &product.DeleteSkuRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewDeleteSkuLogic(context.Background(), ctx)
			resp, err := l.DeleteSku(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
			}
		})
	}
}
