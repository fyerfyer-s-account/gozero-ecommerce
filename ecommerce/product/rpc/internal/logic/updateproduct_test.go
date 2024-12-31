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

func TestUpdateProductLogic_UpdateProduct(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product
	testProduct := &model.Products{
		Name:        "Test Product",
		Description: sql.NullString{String: "Test Description", Valid: true},
		CategoryId:  1,
		Brand:       sql.NullString{String: "Test Brand", Valid: true},
		Price:       99.99,
		Status:      1,
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
		req     *product.UpdateProductRequest
		wantErr error
		verify  func(*testing.T, *model.Products)
	}{
		{
			name: "Update single field",
			req: &product.UpdateProductRequest{
				Id:   productId,
				Name: "Updated Product",
			},
			wantErr: nil,
			verify: func(t *testing.T, p *model.Products) {
				assert.Equal(t, "Updated Product", p.Name)
			},
		},
		{
			name: "Update multiple fields",
			req: &product.UpdateProductRequest{
				Id:          productId,
				Description: "Updated Description",
				Price:       199.99,
				Brand:       "Updated Brand",
			},
			wantErr: nil,
			verify: func(t *testing.T, p *model.Products) {
				assert.Equal(t, "Updated Description", p.Description.String)
				assert.True(t, p.Description.Valid)
				assert.Equal(t, "Updated Brand", p.Brand.String)
				assert.True(t, p.Brand.Valid)
				assert.Equal(t, 199.99, p.Price)
			},
		},
		{
			name: "Invalid ID",
			req: &product.UpdateProductRequest{
				Id:   0,
				Name: "Test",
			},
			wantErr: zeroerr.ErrInvalidParam,
			verify:  nil,
		},
		{
			name: "Empty update",
			req: &product.UpdateProductRequest{
				Id: productId,
			},
			wantErr: zeroerr.ErrInvalidParam,
			verify:  nil,
		},
		{
			name: "Non-existent product",
			req: &product.UpdateProductRequest{
				Id:   99999,
				Name: "Test",
			},
			wantErr: zeroerr.ErrProductNotFound, // Changed from model.ErrNotFound
			verify:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateProductLogic(context.Background(), ctx)
			resp, err := l.UpdateProduct(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				if tt.verify != nil {
					// Verify updates
					updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(productId))
					assert.NoError(t, err)
					tt.verify(t, updated)
				}
			}
		})
	}
}
