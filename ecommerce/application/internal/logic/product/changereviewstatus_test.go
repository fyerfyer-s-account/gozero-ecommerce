package product

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeReviewStatusLogic_ChangeReviewStatus(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.ChangeReviewStatusReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success",
			req: &types.ChangeReviewStatusReq{
				Id:     1,
				Status: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ChangeReviewStatus(
					mock.Anything,
					&product.ChangeReviewStatusRequest{
						Id:     1,
						Status: 1,
					},
				).Return(&product.ChangeReviewStatusResponse{
					Success: true,
				}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Invalid status",
			req: &types.ChangeReviewStatusReq{
				Id:     1,
				Status: 3, // Invalid status value
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ChangeReviewStatus(
					mock.Anything,
					&product.ChangeReviewStatusRequest{
						Id:     1,
						Status: 3,
					},
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Review not found",
			req: &types.ChangeReviewStatusReq{
				Id:     99999,
				Status: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ChangeReviewStatus(
					mock.Anything,
					&product.ChangeReviewStatusRequest{
						Id:     99999,
						Status: 1,
					},
				).Return(nil, zeroerr.ErrReviewNotFound)
			},
			wantErr: zeroerr.ErrReviewNotFound,
		},
		{
			name: "RPC error",
			req: &types.ChangeReviewStatusReq{
				Id:     1,
				Status: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ChangeReviewStatus(
					mock.Anything,
					&product.ChangeReviewStatusRequest{
						Id:     1,
						Status: 1,
					},
				).Return(nil, errors.New("rpc error"))
			},
			wantErr: errors.New("rpc error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			// Create service context with mock
			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			// Create logic instance
			logic := NewChangeReviewStatusLogic(context.Background(), svcCtx)

			// Execute function
			err := logic.ChangeReviewStatus(tt.req)

			// Assert results
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
