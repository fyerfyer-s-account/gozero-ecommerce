package logic

import (
	"context"
	"database/sql"
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

func TestChangeReviewStatusLogic_ChangeReviewStatus(t *testing.T) {
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
	review := &model.ProductReviews{
		ProductId: uint64(productId),
		UserId:    1,
		OrderId:   1,
		Rating:    4,
		Content:   sql.NullString{String: "Test Review", Valid: true},
		Status:    0,
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
		name       string
		req        *product.ChangeReviewStatusRequest
		wantErr    error
		wantStatus int64
	}{
		{
			name: "Change to approved",
			req: &product.ChangeReviewStatusRequest{
				Id:     reviewId,
				Status: 1,
			},
			wantErr:    nil,
			wantStatus: 1,
		},
		{
			name: "Change to rejected",
			req: &product.ChangeReviewStatusRequest{
				Id:     reviewId,
				Status: 2,
			},
			wantErr:    nil,
			wantStatus: 2,
		},
		{
			name: "Invalid status",
			req: &product.ChangeReviewStatusRequest{
				Id:     reviewId,
				Status: 3,
			},
			wantErr:    zeroerr.ErrInvalidParam,
			wantStatus: 2,
		},
		{
			name: "Non-existent review",
			req: &product.ChangeReviewStatusRequest{
				Id:     99999,
				Status: 1,
			},
			wantErr:    zeroerr.ErrReviewNotFound,
			wantStatus: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewChangeReviewStatusLogic(context.Background(), ctx)
			resp, err := l.ChangeReviewStatus(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify status change
				updated, err := ctx.ProductReviewsModel.FindOne(context.Background(), uint64(reviewId))
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, updated.Status)
			}
		})
	}
}
