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

func TestDeleteReviewLogic_DeleteReview(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.DeleteReviewReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success",
			req: &types.DeleteReviewReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteReview(
					mock.Anything,
					&product.DeleteReviewRequest{
						Id: 1,
					},
				).Return(&product.DeleteReviewResponse{
					Success: true,
				}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Review not found",
			req: &types.DeleteReviewReq{
				Id: 999,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteReview(
					mock.Anything,
					&product.DeleteReviewRequest{
						Id: 999,
					},
				).Return(nil, zeroerr.ErrReviewNotFound)
			},
			wantErr: zeroerr.ErrReviewNotFound,
		},
		{
			name: "Invalid review id",
			req: &types.DeleteReviewReq{
				Id: 0,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteReview(
					mock.Anything,
					&product.DeleteReviewRequest{
						Id: 0,
					},
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "No permission",
			req: &types.DeleteReviewReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteReview(
					mock.Anything,
					&product.DeleteReviewRequest{
						Id: 1,
					},
				).Return(nil, zeroerr.ErrNoPermission)
			},
			wantErr: zeroerr.ErrNoPermission,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewDeleteReviewLogic(context.Background(), svcCtx)
			err := logic.DeleteReview(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
