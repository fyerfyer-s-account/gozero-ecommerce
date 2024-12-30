package logic

import (
	"context"
	"database/sql"
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

func TestCreateSkuLogic_CreateSku(t *testing.T) {
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

	defer func() {
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name    string
		req     *product.CreateSkuRequest
		wantErr error
	}{
		{
			name: "Valid SKU",
			req: &product.CreateSkuRequest{
				ProductId: productId,
				SkuCode:   "TEST-SKU-001",
				Price:     99.99,
				Stock:     100,
				Attributes: []*product.SkuAttribute{
					{Key: "color", Value: "red"},
					{Key: "size", Value: "M"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Invalid product ID",
			req: &product.CreateSkuRequest{
				ProductId: 99999,
				SkuCode:   "TEST-SKU-002",
				Price:     99.99,
				Stock:     100,
			},
			wantErr: zeroerr.ErrProductNotFound,
		},
		{
			name: "Invalid price",
			req: &product.CreateSkuRequest{
				ProductId: productId,
				SkuCode:   "TEST-SKU-003",
				Price:     0,
				Stock:     100,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Empty SKU code",
			req: &product.CreateSkuRequest{
				ProductId: productId,
				SkuCode:   "",
				Price:     99.99,
				Stock:     100,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	createdSkus := make([]int64, 0)
	defer func() {
		for _, id := range createdSkus {
			_ = ctx.SkusModel.Delete(context.Background(), uint64(id))
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewCreateSkuLogic(context.Background(), ctx)
			resp, err := l.CreateSku(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Greater(t, resp.Id, int64(0))
				createdSkus = append(createdSkus, resp.Id)
			}
		})
	}
}
