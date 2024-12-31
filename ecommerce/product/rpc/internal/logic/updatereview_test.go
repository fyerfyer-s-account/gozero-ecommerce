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

func TestUpdateReviewLogic_UpdateReview(t *testing.T) {
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

	// Create test review
	images := []string{"test1.jpg", "test2.jpg"}
	imagesJSON, _ := json.Marshal(images)

	review := &model.ProductReviews{
		ProductId: uint64(productId),
		UserId:    1,
		OrderId:   1,
		Rating:    4,
		Content:   sql.NullString{String: "Original review", Valid: true},
		Images:    sql.NullString{String: string(imagesJSON), Valid: true},
		Status:    1,
	}

	result, err = ctx.ProductReviewsModel.Insert(context.Background(), review)
	assert.NoError(t, err)
	reviewId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductReviewsModel.Delete(context.Background(), uint64(reviewId))
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name    string
		req     *product.UpdateReviewRequest
		wantErr error
	}{
		{
			name: "Valid update",
			req: &product.UpdateReviewRequest{
				Id:      reviewId,
				Rating:  5,
				Content: "Updated review content",
				Images:  []string{"new1.jpg", "new2.jpg"},
			},
			wantErr: nil,
		},
		{
			name: "Invalid rating",
			req: &product.UpdateReviewRequest{
				Id:     reviewId,
				Rating: 6,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Invalid review ID",
			req: &product.UpdateReviewRequest{
				Id:     0,
				Rating: 5,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent review",
			req: &product.UpdateReviewRequest{
				Id:     99999,
				Rating: 5,
			},
			wantErr: zeroerr.ErrReviewNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateReviewLogic(context.Background(), ctx)
			resp, err := l.UpdateReview(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify changes
				updated, err := ctx.ProductReviewsModel.FindOne(context.Background(), uint64(reviewId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Rating, int32(updated.Rating))
				assert.Equal(t, tt.req.Content, updated.Content.String)

				var updatedImages []string
				err = json.Unmarshal([]byte(updated.Images.String), &updatedImages)
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Images, updatedImages)

				// Status should remain unchanged
				assert.Equal(t, int64(1), updated.Status)
			}
		})
	}
}
