package logic

import (
	"context"
	"database/sql"
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

func TestUpdateProductSalesLogic_UpdateProductSales(t *testing.T) {
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
		Sales:       100,
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
		req     *product.UpdateProductSalesRequest
		wantErr error
	}{
		{
			name: "Increment sales",
			req: &product.UpdateProductSalesRequest{
				Id:        productId,
				Increment: 10,
			},
			wantErr: nil,
		},
		{
			name: "Decrement sales",
			req: &product.UpdateProductSalesRequest{
				Id:        productId,
				Increment: -5,
			},
			wantErr: nil,
		},
		{
			name: "Invalid ID",
			req: &product.UpdateProductSalesRequest{
				Id:        0,
				Increment: 10,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateProductSalesLogic(context.Background(), ctx)
			resp, err := l.UpdateProductSales(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify sales update
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(productId))
				assert.NoError(t, err)
				if tt.name == "Increment sales" {
					assert.Equal(t, int64(110), updated.Sales)
				} else if tt.name == "Decrement sales" {
					assert.Equal(t, int64(105), updated.Sales)
				}
			}
		})
	}
}
