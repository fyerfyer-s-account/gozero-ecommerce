package logic

import (
	"context"
	"database/sql"
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

func TestUpdateReviewLogic_UpdateReview(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test review
	testReview := &model.ProductReviews{
		ProductId: 1,
		UserId:    1,
		Rating:    4,
		Content:   sql.NullString{String: "Original review", Valid: true},
		Images:    sql.NullString{String: "[]", Valid: true},
		Status:    1,
	}

	result, err := ctx.ProductReviewsModel.Insert(context.Background(), testReview)
	assert.NoError(t, err)
	reviewId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductReviewsModel.Delete(context.Background(), uint64(reviewId))
	}()

	tests := []struct {
		name    string
		req     *product.UpdateReviewRequest
		wantErr bool
		verify  func(*testing.T, uint64, error)
	}{
		{
			name: "Update all fields",
			req: &product.UpdateReviewRequest{
				Id:      reviewId,
				Rating:  5,
				Content: "Updated review",
				Images:  []string{"new_image.jpg"},
				Status:  2,
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				updated, err := ctx.ProductReviewsModel.FindOne(context.Background(), uint64(id))
				assert.NoError(t, err)
				assert.Equal(t, int64(5), updated.Rating)
				assert.Equal(t, "Updated review", updated.Content.String)
				assert.Equal(t, `["new_image.jpg"]`, updated.Images.String)
				assert.Equal(t, int64(2), updated.Status)
			},
		},
		{
			name: "Update partial fields",
			req: &product.UpdateReviewRequest{
				Id:      reviewId,
				Rating:  3,
				Content: "Partially updated",
			},
			wantErr: false,
			verify: func(t *testing.T, id uint64, err error) {
				assert.NoError(t, err)
				updated, err := ctx.ProductReviewsModel.FindOne(context.Background(), uint64(id))
				assert.NoError(t, err)
				assert.Equal(t, int64(3), updated.Rating)
				assert.Equal(t, "Partially updated", updated.Content.String)
			},
		},
		{
			name: "Invalid review ID",
			req: &product.UpdateReviewRequest{
				Id:      0,
				Content: "Invalid update",
			},
			wantErr: true,
			verify: func(t *testing.T, id uint64, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateReviewLogic(context.Background(), ctx)
			resp, err := l.UpdateReview(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
				time.Sleep(100 * time.Millisecond)
				if tt.verify != nil {
					tt.verify(t, uint64(tt.req.Id), err)
				}
			}
		})
	}
}
