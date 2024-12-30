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

func TestCreateCategoryLogic_CreateCategory(t *testing.T) {
	// Load config
	configFile := flag.String("f", "../../etc/product.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test parent category
	parentCategory := &model.Categories{
		Name:      "Test Parent Category",
		ParentId:  0,
		Level:     1,
		Sort:      1,
		Icon:      sql.NullString{String: "test-icon.png", Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := ctx.CategoriesModel.Insert(context.Background(), parentCategory)
	assert.NoError(t, err)
	parentId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup function
	defer func() {
		err := ctx.CategoriesModel.Delete(context.Background(), uint64(parentId))
		assert.NoError(t, err)
	}()

	tests := []struct {
		name    string
		req     *product.CreateCategoryRequest
		wantErr error
	}{
		{
			name: "Valid category",
			req: &product.CreateCategoryRequest{
				Name:     "Test Category",
				ParentId: 0,
				Sort:     1,
				Icon:     "test-icon.png",
			},
			wantErr: nil,
		},
		{
			name: "Valid subcategory",
			req: &product.CreateCategoryRequest{
				Name:     "Test Subcategory",
				ParentId: parentId,
				Sort:     1,
				Icon:     "test-icon.png",
			},
			wantErr: nil,
		},
		{
			name: "Empty name",
			req: &product.CreateCategoryRequest{
				Name:     "",
				ParentId: 0,
				Sort:     1,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Invalid parent",
			req: &product.CreateCategoryRequest{
				Name:     "Test Category",
				ParentId: 99999,
				Sort:     1,
			},
			wantErr: zeroerr.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewCreateCategoryLogic(context.Background(), ctx)
			resp, err := l.CreateCategory(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Greater(t, resp.Id, int64(0))

				// Cleanup created category
				err = ctx.CategoriesModel.Delete(context.Background(), uint64(resp.Id))
				assert.NoError(t, err)
			}
		})
	}
}
