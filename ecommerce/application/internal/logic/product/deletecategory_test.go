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

func TestDeleteCategoryLogic_DeleteCategory(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.DeleteCategoryReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success",
			req: &types.DeleteCategoryReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteCategory(
					mock.Anything,
					&product.DeleteCategoryRequest{
						Id: 1,
					},
				).Return(&product.DeleteCategoryResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Category has children",
			req: &types.DeleteCategoryReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteCategory(
					mock.Anything,
					&product.DeleteCategoryRequest{
						Id: 1,
					},
				).Return(nil, zeroerr.ErrCategoryHasChildren)
			},
			wantErr: zeroerr.ErrCategoryHasChildren,
		},
		{
			name: "Category has products",
			req: &types.DeleteCategoryReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteCategory(
					mock.Anything,
					&product.DeleteCategoryRequest{
						Id: 1,
					},
				).Return(nil, zeroerr.ErrCategoryHasProducts)
			},
			wantErr: zeroerr.ErrCategoryHasProducts,
		},
		{
			name: "Category not found",
			req: &types.DeleteCategoryReq{
				Id: 999,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().DeleteCategory(
					mock.Anything,
					&product.DeleteCategoryRequest{
						Id: 999,
					},
				).Return(nil, zeroerr.ErrCategoryNotFound)
			},
			wantErr: zeroerr.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewDeleteCategoryLogic(context.Background(), svcCtx)
			err := logic.DeleteCategory(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
