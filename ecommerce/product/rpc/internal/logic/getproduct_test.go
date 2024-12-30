package logic

import (
	"context"
	"database/sql"
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

func TestGetProductLogic_GetProduct(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/product.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Prepare images JSON
	images := []string{"image1.jpg", "image2.jpg"}
	imagesJSON, err := json.Marshal(images)
	assert.NoError(t, err)

	// Create test product
	testProduct := &model.Products{
		Name:        "Test Product",
		Description: sql.NullString{String: "Test Description", Valid: true},
		CategoryId:  1,
		Brand:       sql.NullString{String: "Test Brand", Valid: true},
		Images:      sql.NullString{String: string(imagesJSON), Valid: true},
		Price:       99.99,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := ctx.ProductsModel.Insert(context.Background(), testProduct)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test SKU
	skuAttrs := []map[string]string{
		{"key": "color", "value": "red"},
		{"key": "size", "value": "M"},
	}
	skuAttrsJSON, err := json.Marshal(skuAttrs)
	assert.NoError(t, err)

	testSku := &model.Skus{
		ProductId:  uint64(productId),
		SkuCode:    "TEST-SKU-001",
		Attributes: string(skuAttrsJSON),
		Price:      99.99,
		Stock:      100,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	result, err = ctx.SkusModel.Insert(context.Background(), testSku)
	assert.NoError(t, err)
	skuId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		err := ctx.ProductsModel.Delete(context.Background(), uint64(productId))
		assert.NoError(t, err)
		err = ctx.SkusModel.Delete(context.Background(), uint64(skuId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *product.GetProductRequest
		wantErr error
	}{
		{
			name: "Valid product",
			req: &product.GetProductRequest{
				Id: productId,
			},
			wantErr: nil,
		},
		{
			name: "Invalid product ID",
			req: &product.GetProductRequest{
				Id: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent product",
			req: &product.GetProductRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetProductLogic(context.Background(), ctx)
			resp, err := l.GetProduct(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Product)
				assert.Equal(t, testProduct.Name, resp.Product.Name)
				assert.Len(t, resp.Skus, 1)
				assert.Equal(t, testSku.SkuCode, resp.Skus[0].SkuCode)
			}
		})
	}
}
