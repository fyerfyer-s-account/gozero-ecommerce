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

func TestUpdateCategoryLogic_UpdateCategory(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.UpdateCategoryReq
		mock    func(mockProduct *ProductService)
		wantErr error
	}{
		{
			name: "Success - Update name only",
			req: &types.UpdateCategoryReq{
				Id:   1,
				Name: "Updated Category",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateCategory(
					mock.Anything,
					&product.UpdateCategoryRequest{
						Id:   1,
						Name: "Updated Category",
					},
				).Return(&product.UpdateCategoryResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Success - Update all fields",
			req: &types.UpdateCategoryReq{
				Id:   1,
				Name: "Updated Category",
				Sort: 2,
				Icon: "new-icon.png",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateCategory(
					mock.Anything,
					&product.UpdateCategoryRequest{
						Id:   1,
						Name: "Updated Category",
						Sort: 2,
						Icon: "new-icon.png",
					},
				).Return(&product.UpdateCategoryResponse{Success: true}, nil)
			},
			wantErr: nil,
		},
		{
			name: "Category not found",
			req: &types.UpdateCategoryReq{
				Id:   999,
				Name: "Updated Category",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateCategory(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrCategoryNotFound)
			},
			wantErr: zeroerr.ErrCategoryNotFound,
		},
		{
			name: "Duplicate name",
			req: &types.UpdateCategoryReq{
				Id:   1,
				Name: "Existing Category",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().UpdateCategory(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrCategoryDuplicate)
			},
			wantErr: zeroerr.ErrCategoryDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewUpdateCategoryLogic(context.Background(), svcCtx)
			err := logic.UpdateCategory(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
