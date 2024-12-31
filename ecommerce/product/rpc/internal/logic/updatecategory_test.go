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

func TestUpdateCategoryLogic_UpdateCategory(t *testing.T) {
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

	// Create another category for duplicate test
	other := &model.Categories{
		Name:  "Other Category",
		Level: 1,
		Sort:  2,
	}
	_, err = ctx.CategoriesModel.Insert(context.Background(), other)
	assert.NoError(t, err)

	defer func() {
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(categoryId))
	}()

	tests := []struct {
		name    string
		req     *product.UpdateCategoryRequest
		wantErr error
	}{
		{
			name: "Valid update",
			req: &product.UpdateCategoryRequest{
				Id:   categoryId,
				Name: "Updated Category",
				Sort: 2,
				Icon: "new-icon.png",
			},
			wantErr: nil,
		},
		{
			name: "Invalid ID",
			req: &product.UpdateCategoryRequest{
				Id:   0,
				Name: "Test",
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Empty name",
			req: &product.UpdateCategoryRequest{
				Id:   categoryId,
				Name: "",
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Duplicate name",
			req: &product.UpdateCategoryRequest{
				Id:   categoryId,
				Name: "Other Category",
			},
			wantErr: zeroerr.ErrCategoryDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateCategoryLogic(context.Background(), ctx)
			resp, err := l.UpdateCategory(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify changes
				updated, err := ctx.CategoriesModel.FindOne(context.Background(), uint64(categoryId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Name, updated.Name)
				assert.Equal(t, tt.req.Sort, updated.Sort)
				assert.Equal(t, tt.req.Icon, updated.Icon.String)
			}
		})
	}
}
