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

func TestUpdateProductLogic_UpdateProduct(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.UpdateProductReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success - Update name only",
			req: &types.UpdateProductReq{
				Id:   1,
				Name: "Updated Product",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProduct(
					mock.Anything,
					mock.MatchedBy(func(req *product.UpdateProductRequest) bool {
						return req.Id == 1 && req.Name == "Updated Product"
					}),
				).Return(&product.UpdateProductResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Success - Update multiple fields",
			req: &types.UpdateProductReq{
				Id:          1,
				Name:        "Updated Product",
				Description: "New Description",
				Price:       199.99,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProduct(
					mock.Anything,
					mock.MatchedBy(func(req *product.UpdateProductRequest) bool {
						return req.Id == 1 &&
							req.Name == "Updated Product" &&
							req.Description == "New Description" &&
							req.Price == 199.99
					}),
				).Return(&product.UpdateProductResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Empty update",
			req: &types.UpdateProductReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProduct(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Product not found",
			req: &types.UpdateProductReq{
				Id:   999,
				Name: "Updated Product",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateProduct(
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

			logic := NewUpdateProductLogic(context.Background(), svcCtx)
			err := logic.UpdateProduct(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
