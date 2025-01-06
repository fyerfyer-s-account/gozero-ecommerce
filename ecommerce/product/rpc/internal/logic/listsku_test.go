package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestListSkusLogic_ListSkus(t *testing.T) {
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

	// Create test SKUs
	attrs := []map[string]string{
		{"color": "red", "size": "M"},
		{"color": "blue", "size": "L"},
	}
	attrsJSON, _ := json.Marshal(attrs)

	skus := []*model.Skus{
		{
			ProductId:  uint64(productId),
			SkuCode:    "SKU-001",
			Price:      99.99,
			Stock:      100,
			Attributes: sql.NullString {
				String: string(attrsJSON),
				Valid:  true,
			},
		},
		{
			ProductId:  uint64(productId),
			SkuCode:    "SKU-002",
			Price:      109.99,
			Stock:      50,
			Attributes: sql.NullString {
				String: string(attrsJSON),
				Valid:  true,
			},
		},
	}

	skuIds := make([]int64, 0)
	for _, sku := range skus {
		result, err := ctx.SkusModel.Insert(context.Background(), sku)
		assert.NoError(t, err)
		id, err := result.LastInsertId()
		assert.NoError(t, err)
		skuIds = append(skuIds, id)
	}

	defer func() {
		for _, id := range skuIds {
			_ = ctx.SkusModel.Delete(context.Background(), uint64(id))
		}
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name      string
		req       *product.ListSkusRequest
		wantCount int
		wantErr   bool
	}{
		{
			name: "List all SKUs",
			req: &product.ListSkusRequest{
				ProductId: productId,
				Page:      1,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Test pagination",
			req: &product.ListSkusRequest{
				ProductId: productId,
				Page:      2,
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "Invalid product ID",
			req: &product.ListSkusRequest{
				ProductId: 0,
				Page:      1,
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListSkusLogic(context.Background(), ctx)
			resp, err := l.ListSkus(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantCount, int64(len(resp.Skus)))

				if tt.wantCount > 0 {
					for _, s := range resp.Skus {
						assert.Equal(t, productId, s.ProductId)
						assert.NotEmpty(t, s.SkuCode)
						assert.Greater(t, s.Price, float64(0))
						assert.Greater(t, s.Stock, int64(0))
						assert.NotEmpty(t, s.Attributes)
					}
				}
			}
		})
	}
}
