package product

import (
	"context"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateReviewLogic_UpdateReview(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.UpdateReviewReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success - Update all fields",
			req: &types.UpdateReviewReq{
				Id:      1,
				Rating:  5,
				Content: "Updated review content",
				Images:  []string{"new1.jpg", "new2.jpg"},
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateReview(
					mock.Anything,
					&product.UpdateReviewRequest{
						Id:      1,
						Rating:  5,
						Content: "Updated review content",
						Images:  []string{"new1.jpg", "new2.jpg"},
					},
				).Return(&product.UpdateReviewResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Success - Update rating only",
			req: &types.UpdateReviewReq{
				Id:     1,
				Rating: 4,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateReview(
					mock.Anything,
					&product.UpdateReviewRequest{
						Id:     1,
						Rating: 4,
					},
				).Return(&product.UpdateReviewResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Review not found",
			req: &types.UpdateReviewReq{
				Id:      999,
				Rating:  5,
				Content: "Test content",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateReview(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrReviewNotFound)
			},
			wantErr: zeroerr.ErrReviewNotFound,
		},
		{
			name: "Invalid rating",
			req: &types.UpdateReviewReq{
				Id:     1,
				Rating: 6,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateReview(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewUpdateReviewLogic(context.Background(), svcCtx)
			err := logic.UpdateReview(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
