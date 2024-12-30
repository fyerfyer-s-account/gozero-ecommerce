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

func TestGetCategoryLogic_GetCategory(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/product.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test parent category
	testCategory := &model.Categories{
		Name:      "Test Category",
		ParentId:  0,
		Level:     1,
		Sort:      1,
		Icon:      sql.NullString{String: "test-icon.png", Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := ctx.CategoriesModel.Insert(context.Background(), testCategory)
	assert.NoError(t, err)
	categoryId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test child category
	testChildCategory := &model.Categories{
		Name:      "Test Child Category",
		ParentId:  uint64(categoryId),
		Level:     2,
		Sort:      1,
		Icon:      sql.NullString{String: "test-child-icon.png", Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err = ctx.CategoriesModel.Insert(context.Background(), testChildCategory)
	assert.NoError(t, err)
	childId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		err := ctx.CategoriesModel.Delete(context.Background(), uint64(childId))
		assert.NoError(t, err)
		err = ctx.CategoriesModel.Delete(context.Background(), uint64(categoryId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *product.GetCategoryRequest
		wantErr error
	}{
		{
			name: "Valid category",
			req: &product.GetCategoryRequest{
				Id: categoryId,
			},
			wantErr: nil,
		},
		{
			name: "Invalid category ID",
			req: &product.GetCategoryRequest{
				Id: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent category",
			req: &product.GetCategoryRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetCategoryLogic(context.Background(), ctx)
			resp, err := l.GetCategory(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Category)
				assert.Equal(t, testCategory.Name, resp.Category.Name)
				assert.Len(t, resp.Children, 1)
				assert.Equal(t, testChildCategory.Name, resp.Children[0].Name)
			}
		})
	}
}
