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

func TestCreateReviewLogic_CreateReview(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product
	testProduct := &model.Products{
		Name:        "Test Product",
		Description: sql.NullString{String: "Test Description", Valid: true},
		CategoryId:  1,
		Price:       99.99,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := ctx.ProductsModel.Insert(context.Background(), testProduct)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name    string
		req     *product.CreateReviewRequest
		wantErr error
	}{
		{
			name: "Valid review",
			req: &product.CreateReviewRequest{
				UserId:    1,
				ProductId: productId,
				OrderId:   1,
				Rating:    5,
				Content:   "Great product, highly recommended!",
				Images:    []string{"review1.jpg", "review2.jpg"},
			},
			wantErr: nil,
		},
		{
			name: "Invalid product ID",
			req: &product.CreateReviewRequest{
				UserId:    1,
				ProductId: 99999,
				OrderId:   1,
				Rating:    5,
				Content:   "Great product!",
			},
			wantErr: zeroerr.ErrProductNotFound,
		},
		{
			name: "Invalid rating",
			req: &product.CreateReviewRequest{
				UserId:    1,
				ProductId: productId,
				OrderId:   1,
				Rating:    6,
				Content:   "Great product!",
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Content too short",
			req: &product.CreateReviewRequest{
				UserId:    1,
				ProductId: productId,
				OrderId:   1,
				Rating:    5,
				Content:   "OK",
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	createdReviews := make([]int64, 0)
	defer func() {
		for _, id := range createdReviews {
			_ = ctx.ProductReviewsModel.Delete(context.Background(), uint64(id))
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewCreateReviewLogic(context.Background(), ctx)
			resp, err := l.CreateReview(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Greater(t, resp.Id, int64(0))
				createdReviews = append(createdReviews, resp.Id)
			}
		})
	}
}
