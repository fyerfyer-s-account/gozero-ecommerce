package product

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProductReviewsLogic_GetProductReviews(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.ReviewListReq
		mock    func(mockProduct *ProductService)
		want    []types.Review
		wantErr error
	}{
		{
			name: "successful get reviews",
			req: &types.ReviewListReq{
				ProductId: 1,
				Page:      1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListReviews(
					mock.Anything,
					&product.ListReviewsRequest{
						ProductId: 1,
						Page:      1,
					},
				).Return(&product.ListReviewsResponse{
					Total: 1,
					Reviews: []*product.Review{
						{
							Id:        1,
							ProductId: 1,
							OrderId:   100,
							UserId:    1000,
							Rating:    5,
							Content:   "Great product!",
							Images:    []string{"image1.jpg", "image2.jpg"},
							CreatedAt: 1234567890,
						},
					},
				}, nil)
			},
			want: []types.Review{
				{
					Id:        1,
					ProductId: 1,
					OrderId:   100,
					UserId:    1000,
					Rating:    5,
					Content:   "Great product!",
					Images:    []string{"image1.jpg", "image2.jpg"},
					CreatedAt: 1234567890,
				},
			},
			wantErr: nil,
		},
		{
			name: "rpc error",
			req: &types.ReviewListReq{
				ProductId: 1,
				Page:      1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListReviews(
					mock.Anything,
					&product.ListReviewsRequest{
						ProductId: 1,
						Page:      1,
					},
				).Return(nil, errors.New("rpc error"))
			},
			want:    nil,
			wantErr: errors.New("rpc error"),
		},
		{
			name: "empty reviews",
			req: &types.ReviewListReq{
				ProductId: 1,
				Page:      1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListReviews(
					mock.Anything,
					&product.ListReviewsRequest{
						ProductId: 1,
						Page:      1,
					},
				).Return(&product.ListReviewsResponse{
					Total:   0,
					Reviews: []*product.Review{},
				}, nil)
			},
			want:    []types.Review{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewGetProductReviewsLogic(context.Background(), svcCtx)
			got, err := logic.GetProductReviews(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
