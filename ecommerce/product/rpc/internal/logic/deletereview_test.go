package logic

import (
	"context"
	"database/sql"
	"encoding/json"
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

func TestDeleteReviewLogic_DeleteReview(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product first
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
	images := []string{"review1.jpg", "review2.jpg"}
	imagesJSON, err := json.Marshal(images)
	assert.NoError(t, err)

	testReview := &model.ProductReviews{
		ProductId: uint64(productId),
		UserId:    1,
		OrderId:   1,
		Rating:    5,
		Content:   sql.NullString{String: "Test Review", Valid: true},
		Images:    sql.NullString{String: string(imagesJSON), Valid: true},
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err = ctx.ProductReviewsModel.Insert(context.Background(), testReview)
	assert.NoError(t, err)
	reviewId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name    string
		req     *product.DeleteReviewRequest
		wantErr error
	}{
		{
			name: "Valid review",
			req: &product.DeleteReviewRequest{
				Id: reviewId,
			},
			wantErr: nil,
		},
		{
			name: "Invalid review ID",
			req: &product.DeleteReviewRequest{
				Id: 0,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent review",
			req: &product.DeleteReviewRequest{
				Id: 99999,
			},
			wantErr: zeroerr.ErrReviewNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewDeleteReviewLogic(context.Background(), ctx)
			resp, err := l.DeleteReview(tt.req)

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
