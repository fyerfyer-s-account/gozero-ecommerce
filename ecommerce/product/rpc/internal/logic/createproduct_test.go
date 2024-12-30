package logic

import (
	"context"
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

func TestCreateProductLogic_CreateProduct(t *testing.T) {
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

	tests := []struct {
		name    string
		req     *product.CreateProductRequest
		wantErr error
	}{
		{
			name: "Valid product",
			req: &product.CreateProductRequest{
				Name:        "Test Product",
				Description: "Test Description",
				CategoryId:  categoryId,
				Brand:       "Test Brand",
				Images:      []string{"image1.jpg", "image2.jpg"},
				Price:       99.99,
				SkuAttributes: []*product.SkuAttribute{
					{Key: "color", Value: "red"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Invalid category",
			req: &product.CreateProductRequest{
				Name:       "Test Product 2",
				CategoryId: 99999,
				Price:      99.99,
			},
			wantErr: zeroerr.ErrCategoryNotFound,
		},
		{
			name: "Invalid price",
			req: &product.CreateProductRequest{
				Name:       "Test Product 3",
				CategoryId: categoryId,
				Price:      0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	createdIds := make([]int64, 0)
	defer func() {
		for _, id := range createdIds {
			_ = ctx.ProductsModel.Delete(context.Background(), uint64(id))
			_ = ctx.SkusModel.DeleteByProductId(context.Background(), uint64(id))
		}
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(categoryId))
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewCreateProductLogic(context.Background(), ctx)
			resp, err := l.CreateProduct(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Greater(t, resp.Id, int64(0))
				createdIds = append(createdIds, resp.Id)
			}
		})
	}
}
