package logic

import (
	"context"
	"database/sql"
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

func TestUpdateProductLogic_UpdateProduct(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test category
	category := &model.Categories{
		Name:  "Test Category",
		Level: 1,
		Sort:  1,
	}
	result, err := ctx.CategoriesModel.Insert(context.Background(), category)
	assert.NoError(t, err)
	categoryId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test product
	images := []string{"test1.jpg", "test2.jpg"}
	imagesJSON, _ := json.Marshal(images)

	p := &model.Products{
		Name:        "Test Product",
		Description: sql.NullString{String: "Test Description", Valid: true},
		CategoryId:  uint64(categoryId),
		Brand:       sql.NullString{String: "Test Brand", Valid: true},
		Images:      sql.NullString{String: string(imagesJSON), Valid: true},
		Price:       99.99,
		Status:      1,
	}

	result, err = ctx.ProductsModel.Insert(context.Background(), p)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(categoryId))
	}()

	tests := []struct {
		name    string
		req     *product.UpdateProductRequest
		wantErr error
	}{
		{
			name: "Valid update",
			req: &product.UpdateProductRequest{
				Id:          productId,
				Name:        "Updated Product",
				Description: "Updated Description",
				CategoryId:  categoryId,
				Brand:       "Updated Brand",
				Images:      []string{"new1.jpg", "new2.jpg"},
				Price:       199.99,
				Status:      1,
			},
			wantErr: nil,
		},
		{
			name: "Invalid ID",
			req: &product.UpdateProductRequest{
				Id:    0,
				Name:  "Test",
				Price: 99.99,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Invalid price",
			req: &product.UpdateProductRequest{
				Id:    productId,
				Name:  "Test",
				Price: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent product",
			req: &product.UpdateProductRequest{
				Id:    99999,
				Name:  "Test",
				Price: 99.99,
			},
			wantErr: zeroerr.ErrProductNotFound,
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

				// Verify changes
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(productId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Name, updated.Name)
				assert.Equal(t, tt.req.Description, updated.Description.String)
				assert.Equal(t, tt.req.Brand, updated.Brand.String)
				assert.Equal(t, tt.req.Price, updated.Price)

				var updatedImages []string
				err = json.Unmarshal([]byte(updated.Images.String), &updatedImages)
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Images, updatedImages)
			}
		})
	}
}
