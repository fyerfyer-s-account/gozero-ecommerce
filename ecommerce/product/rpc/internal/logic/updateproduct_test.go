package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"

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
		Description: sql.NullString{String: "Original description", Valid: true},
		CategoryId:  1,
		Brand:       sql.NullString{String: "Original brand", Valid: true},
		Price:       99.99,
		Status:      1,
		Sales:       10,
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
		wantErr bool
		verify  func(*testing.T, uint64, error)
	}{
		{
			name: "Update all fields",
			req: &product.UpdateProductRequest{
				Id:          productId,
				Name:        "Updated Product",
				Description: "Updated description",
				CategoryId:  2,
				Brand:       "Updated brand",
				Price:       199.99,
				Status:      2,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(id))
				assert.NoError(t, err)
				assert.Equal(t, "Updated Product", updated.Name)
				assert.Equal(t, "Updated description", updated.Description.String)
				assert.Equal(t, uint64(2), updated.CategoryId)
				assert.Equal(t, "Updated brand", updated.Brand.String)
				assert.Equal(t, 199.99, updated.Price)
				assert.Equal(t, int64(2), updated.Status)
			},
		},
		{
			name: "Update partial fields",
			req: &product.UpdateProductRequest{
				Id:    productId,
				Name:  "Partially Updated",
				Price: 299.99,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(id))
				assert.NoError(t, err)
				assert.Equal(t, "Partially Updated", updated.Name)
				assert.Equal(t, 299.99, updated.Price)
			},
		},
		{
			name: "Update with sales increment",
			req: &product.UpdateProductRequest{
				Id:             productId,
				SalesIncrement: 5,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(id))
				assert.NoError(t, err)
				assert.Equal(t, int64(15), updated.Sales) // 10 + 5
			},
		},
		{
			name: "Invalid product ID",
			req: &product.UpdateProductRequest{
				Id:   0,
				Name: "Invalid Update",
			},
			wantErr: true,
			verify: func(t *testing.T, id uint64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateProductLogic(context.Background(), ctx)
			resp, err := l.UpdateProduct(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
				if tt.verify != nil {
					tt.verify(t, uint64(tt.req.Id), err)
				}
			}
		})
	}
}
