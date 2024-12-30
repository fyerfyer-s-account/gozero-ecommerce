package logic

import (
	"context"
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

func TestDeleteCategoryLogic_DeleteCategory(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create parent category
	parent := &model.Categories{
		Name:      "Test Parent",
		ParentId:  0,
		Level:     1,
		Sort:      1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err := ctx.CategoriesModel.Insert(context.Background(), parent)
	assert.NoError(t, err)
	parentId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Wait for parent category to be created
	time.Sleep(time.Millisecond * 100)

	// Create child category
	child := &model.Categories{
		Name:      "Test Child",
		ParentId:  uint64(parentId),
		Level:     2,
		Sort:      1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err = ctx.CategoriesModel.Insert(context.Background(), child)
	assert.NoError(t, err)
	childId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Verify parent-child relationship before testing delete
	hasChildren, err := ctx.CategoriesModel.HasChildren(context.Background(), uint64(parentId))
	assert.NoError(t, err)
	assert.True(t, hasChildren, "Parent category should have children")

	// Create category with product
	withProduct := &model.Categories{
		Name:      "Test With Product",
		ParentId:  0,
		Level:     1,
		Sort:      2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result, err = ctx.CategoriesModel.Insert(context.Background(), withProduct)
	assert.NoError(t, err)
	withProductId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create product under category
	p := &model.Products{
		Name:       "Test Product",
		CategoryId: uint64(withProductId),
		Price:      99.99,
		Status:     1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err = ctx.ProductsModel.Insert(context.Background(), p)
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(childId))
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(parentId))
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(withProductId))
	}()

	tests := []struct {
		name    string
		req     *product.DeleteCategoryRequest
		wantErr error
	}{
		{
			name: "Delete leaf category",
			req: &product.DeleteCategoryRequest{
				Id: childId,
			},
			wantErr: nil,
		},
		{
			name: "Delete category with children",
			req: &product.DeleteCategoryRequest{
				Id: parentId,
			},
			wantErr: zeroerr.ErrCategoryHasChildren,
		},
		{
			name: "Delete category with products",
			req: &product.DeleteCategoryRequest{
				Id: withProductId,
			},
			wantErr: zeroerr.ErrCategoryHasProducts,
		},
		{
			name: "Delete non-existent category",
			req: &product.DeleteCategoryRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewDeleteCategoryLogic(context.Background(), ctx)
			resp, err := l.DeleteCategory(tt.req)

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
