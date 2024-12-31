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

func TestUpdateProductStatusLogic_UpdateProductStatus(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.UpdateProductStatusReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success - Online status",
			req: &types.UpdateProductStatusReq{
				Id:     1,
				Status: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProductStatus(
					mock.Anything,
					&product.UpdateProductStatusRequest{
						Id:     1,
						Status: 1,
					},
				).Return(&product.UpdateProductStatusResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Success - Offline status",
			req: &types.UpdateProductStatusReq{
				Id:     1,
				Status: 2,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProductStatus(
					mock.Anything,
					&product.UpdateProductStatusRequest{
						Id:     1,
						Status: 2,
					},
				).Return(&product.UpdateProductStatusResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Invalid status value",
			req: &types.UpdateProductStatusReq{
				Id:     1,
				Status: 3,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProductStatus(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Product not found",
			req: &types.UpdateProductStatusReq{
				Id:     999,
				Status: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProductStatus(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrProductNotFound)
			},
			wantErr: zeroerr.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewUpdateProductStatusLogic(context.Background(), svcCtx)
			err := logic.UpdateProductStatus(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
