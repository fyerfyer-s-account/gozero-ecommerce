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

func TestUpdateSkuPriceLogic_UpdateSkuPrice(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.UpdateSkuPriceReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success",
			req: &types.UpdateSkuPriceReq{
				Id:    1,
				Price: 199.99,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuPrice(
					mock.Anything,
					&product.UpdateSkuPriceRequest{
						Id:    1,
						Price: 199.99,
					},
				).Return(&product.UpdateSkuPriceResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Invalid price",
			req: &types.UpdateSkuPriceReq{
				Id:    1,
				Price: 0,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuPrice(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrInvalidProductPrice)
			},
			wantErr: zeroerr.ErrInvalidProductPrice,
		},
		{
			name: "SKU not found",
			req: &types.UpdateSkuPriceReq{
				Id:    999,
				Price: 199.99,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateSkuPrice(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrSkuNotFound)
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewUpdateSkuPriceLogic(context.Background(), svcCtx)
			err := logic.UpdateSkuPrice(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}