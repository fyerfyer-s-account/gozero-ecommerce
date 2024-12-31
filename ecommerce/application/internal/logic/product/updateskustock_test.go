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

func TestUpdateSkuStockLogic_UpdateSkuStock(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.UpdateSkuStockReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success - Increase stock",
			req: &types.UpdateSkuStockReq{
				Id:    1,
				Stock: 100,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuStock(
					mock.Anything,
					&product.UpdateSkuStockRequest{
						Id:        1,
						Increment: 100,
					},
				).Return(&product.UpdateSkuStockResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Success - Decrease stock",
			req: &types.UpdateSkuStockReq{
				Id:    1,
				Stock: -50,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuStock(
					mock.Anything,
					&product.UpdateSkuStockRequest{
						Id:        1,
						Increment: -50,
					},
				).Return(&product.UpdateSkuStockResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "SKU not found",
			req: &types.UpdateSkuStockReq{
				Id:    999,
				Stock: 100,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuStock(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrSkuNotFound)
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
		{
			name: "Invalid stock (would result in negative)",
			req: &types.UpdateSkuStockReq{
				Id:    1,
				Stock: -200,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuStock(
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

			logic := NewUpdateSkuStockLogic(context.Background(), svcCtx)
			err := logic.UpdateSkuStock(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
