package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestListProductsLogic_ListProducts(t *testing.T) {
	// Setup
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test categories
	category1 := &model.Categories{
		Name:     "Category 1",
		ParentId: 0,
		Level:    1,
		Sort:     1,
	}
	category2 := &model.Categories{
		Name:     "Category 2",
		ParentId: 0,
		Level:    1,
		Sort:     2,
	}

	result, err := ctx.CategoriesModel.Insert(context.Background(), category1)
	assert.NoError(t, err)
	category1Id, err := result.LastInsertId()
	assert.NoError(t, err)

	result, err = ctx.CategoriesModel.Insert(context.Background(), category2)
	assert.NoError(t, err)
	category2Id, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test products
	images := []string{"test1.jpg", "test2.jpg"}
	imagesJSON, _ := json.Marshal(images)

	products := []*model.Products{
		{
			Name:        "iPhone 13",
			Description: sql.NullString{String: "Apple iPhone 13", Valid: true},
			CategoryId:  uint64(category1Id),
			Brand:       sql.NullString{String: "Apple", Valid: true},
			Images:      sql.NullString{String: string(imagesJSON), Valid: true},
			Price:       999.99,
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Samsung S21",
			Description: sql.NullString{String: "Samsung Galaxy S21", Valid: true},
			CategoryId:  uint64(category1Id),
			Brand:       sql.NullString{String: "Samsung", Valid: true},
			Images:      sql.NullString{String: string(imagesJSON), Valid: true},
			Price:       899.99,
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "iPad Pro",
			Description: sql.NullString{String: "Apple iPad Pro", Valid: true},
			CategoryId:  uint64(category2Id),
			Brand:       sql.NullString{String: "Apple", Valid: true},
			Images:      sql.NullString{String: string(imagesJSON), Valid: true},
			Price:       799.99,
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	productIds := make([]int64, 0)
	for _, p := range products {
		result, err := ctx.ProductsModel.Insert(context.Background(), p)
		assert.NoError(t, err)
		id, err := result.LastInsertId()
		assert.NoError(t, err)
		productIds = append(productIds, id)
	}

	// Cleanup
	defer func() {
		for _, id := range productIds {
			_ = ctx.ProductsModel.Delete(context.Background(), uint64(id))
		}
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(category1Id))
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(category2Id))
	}()

	tests := []struct {
		name      string
		req       *product.ListProductsRequest
		wantCount int64
		wantErr   bool
	}{
		{
			name: "List by category 1",
			req: &product.ListProductsRequest{
				CategoryId: category1Id,
				Page:       1,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "List by category 2",
			req: &product.ListProductsRequest{
				CategoryId: category2Id,
				Page:       1,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Search by keyword 'Apple'",
			req: &product.ListProductsRequest{
				Keyword:  "Apple",
				Page:     1,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "General search",
			req: &product.ListProductsRequest{
				Page:     1,
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "Test pagination",
			req: &product.ListProductsRequest{
				Page:     1,
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListProductsLogic(context.Background(), ctx)
			resp, err := l.ListProducts(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Products, int(tt.wantCount))
				for _, p := range resp.Products {
					assert.NotEmpty(t, p.Name)
					assert.NotEmpty(t, p.Description)
					assert.NotEmpty(t, p.Brand)
					assert.Greater(t, p.Price, float64(0))
					assert.NotEmpty(t, p.Images)
				}
			}
		})
	}
}
