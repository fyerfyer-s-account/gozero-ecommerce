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

func TestDeleteProductLogic_DeleteProduct(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product
	testProduct := &model.Products{
		Name:        "Test Product",
		Description: sql.NullString{String: "Test Description", Valid: true},
		CategoryId:  1,
		Price:       99.99,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := ctx.ProductsModel.Insert(context.Background(), testProduct)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test SKU with valid attributes JSON
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

	_, err = ctx.SkusModel.Insert(context.Background(), testSku)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		req     *product.DeleteProductRequest
		wantErr error
	}{
		{
			name: "Valid product",
			req: &product.DeleteProductRequest{
				Id: productId,
			},
			wantErr: nil,
		},
		{
			name: "Invalid product ID",
			req: &product.DeleteProductRequest{
				Id: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent product",
			req: &product.DeleteProductRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewDeleteProductLogic(context.Background(), ctx)
			resp, err := l.DeleteProduct(tt.req)

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
