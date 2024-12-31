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

func TestDeleteProductLogic_DeleteProduct(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.DeleteProductReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success",
			req: &types.DeleteProductReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteProduct(
					mock.Anything,
					&product.DeleteProductRequest{
						Id: 1,
					},
				).Return(&product.DeleteProductResponse{
					Success: true,
				}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Product not found",
			req: &types.DeleteProductReq{
				Id: 999,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteProduct(
					mock.Anything,
					&product.DeleteProductRequest{
						Id: 999,
					},
				).Return(nil, zeroerr.ErrProductNotFound)
			},
			wantErr: zeroerr.ErrProductNotFound,
		},
		{
			name: "Product has orders",
			req: &types.DeleteProductReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteProduct(
					mock.Anything,
					&product.DeleteProductRequest{
						Id: 1,
					},
				).Return(nil, zeroerr.ErrProductHasOrders)
			},
			wantErr: zeroerr.ErrProductHasOrders,
		},
		{
			name: "Product has reviews",
			req: &types.DeleteProductReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteProduct(
					mock.Anything,
					&product.DeleteProductRequest{
						Id: 1,
					},
				).Return(nil, zeroerr.ErrProductHasReviews)
			},
			wantErr: zeroerr.ErrProductHasReviews,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewDeleteProductLogic(context.Background(), svcCtx)
			err := logic.DeleteProduct(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
