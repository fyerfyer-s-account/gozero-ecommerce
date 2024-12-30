package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestListReviewsLogic_ListReviews(t *testing.T) {
	// Setup
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product
	testProduct := &model.Products{
		Name:   "Test Product",
		Price:  99.99,
		Status: 1,
	}
	result, err := ctx.ProductsModel.Insert(context.Background(), testProduct)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create test reviews
	images := []string{"review1.jpg", "review2.jpg"}
	imagesJSON, _ := json.Marshal(images)

	reviews := []*model.ProductReviews{
		{
			ProductId: uint64(productId),
			UserId:    1,
			OrderId:   1,
			Rating:    5,
			Content:   sql.NullString{String: "Great product!", Valid: true},
			Images:    sql.NullString{String: string(imagesJSON), Valid: true},
			Status:    1,
		},
		{
			ProductId: uint64(productId),
			UserId:    2,
			OrderId:   2,
			Rating:    4,
			Content:   sql.NullString{String: "Good product!", Valid: true},
			Images:    sql.NullString{String: string(imagesJSON), Valid: true},
			Status:    1,
		},
	}

	reviewIds := make([]int64, 0)
	for _, review := range reviews {
		result, err := ctx.ProductReviewsModel.Insert(context.Background(), review)
		assert.NoError(t, err)
		id, err := result.LastInsertId()
		assert.NoError(t, err)
		reviewIds = append(reviewIds, id)
	}

	// Cleanup
	defer func() {
		for _, id := range reviewIds {
			_ = ctx.ProductReviewsModel.Delete(context.Background(), uint64(id))
		}
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name      string
		req       *product.ListReviewsRequest
		wantCount int64
		wantErr   bool
	}{
		{
			name: "List all reviews",
			req: &product.ListReviewsRequest{
				ProductId: productId,
				Page:      1,
				PageSize:  10,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Test pagination",
			req: &product.ListReviewsRequest{
				ProductId: productId,
				Page:      1,
				PageSize:  1,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "Invalid product ID",
			req: &product.ListReviewsRequest{
				ProductId: 0,
				Page:      1,
				PageSize:  10,
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListReviewsLogic(context.Background(), ctx)
			resp, err := l.ListReviews(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantCount, int64(len(resp.Reviews)))

				if tt.wantCount > 0 {
					for _, r := range resp.Reviews {
						assert.Equal(t, productId, r.ProductId)
						assert.NotEmpty(t, r.Content)
						assert.Greater(t, r.Rating, int32(0))
						assert.NotEmpty(t, r.Images)
					}
				}
			}
		})
	}
}
