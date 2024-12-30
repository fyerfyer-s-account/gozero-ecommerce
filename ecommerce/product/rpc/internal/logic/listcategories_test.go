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

func TestListCategoriesLogic_ListCategories(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test categories
	parent := &model.Categories{
		Name:     "Parent Category",
		ParentId: 0,
		Level:    1,
		Sort:     1,
		Icon:     sql.NullString{String: "parent.png", Valid: true},
	}
	result, err := ctx.CategoriesModel.Insert(context.Background(), parent)
	assert.NoError(t, err)
	parentId, err := result.LastInsertId()
	assert.NoError(t, err)

	child1 := &model.Categories{
		Name:     "Child Category 1",
		ParentId: uint64(parentId),
		Level:    2,
		Sort:     1,
		Icon:     sql.NullString{String: "child1.png", Valid: true},
	}
	result, err = ctx.CategoriesModel.Insert(context.Background(), child1)
	assert.NoError(t, err)
	child1Id, err := result.LastInsertId()
	assert.NoError(t, err)

	child2 := &model.Categories{
		Name:     "Child Category 2",
		ParentId: uint64(parentId),
		Level:    2,
		Sort:     2,
		Icon:     sql.NullString{String: "child2.png", Valid: true},
	}
	result, err = ctx.CategoriesModel.Insert(context.Background(), child2)
	assert.NoError(t, err)
	child2Id, err := result.LastInsertId()
	assert.NoError(t, err)

	// Cleanup
	defer func() {
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(child1Id))
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(child2Id))
		_ = ctx.CategoriesModel.Delete(context.Background(), uint64(parentId))
	}()

	tests := []struct {
		name         string
		req          *product.ListCategoriesRequest
		wantCount    int
		wantParentId int64
	}{
		{
			name: "List root categories",
			req: &product.ListCategoriesRequest{
				ParentId: 0,
			},
			wantCount:    1,
			wantParentId: 0,
		},
		{
			name: "List child categories",
			req: &product.ListCategoriesRequest{
				ParentId: parentId,
			},
			wantCount:    2,
			wantParentId: parentId,
		},
		{
			name: "List non-existent parent",
			req: &product.ListCategoriesRequest{
				ParentId: 99999,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListCategoriesLogic(context.Background(), ctx)
			resp, err := l.ListCategories(tt.req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp.Categories, tt.wantCount)

			if tt.wantCount > 0 {
				for _, cat := range resp.Categories {
					assert.Equal(t, tt.wantParentId, cat.ParentId)
				}
			}
		})
	}
}
