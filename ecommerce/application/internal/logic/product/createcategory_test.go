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

func TestCreateCategoryLogic_CreateCategory(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.CreateCategoryReq
		mock    func(mockProduct *ProductService)
		want    *types.CreateCategoryResp
		wantErr error
	}{
		{
			name: "Success",
			req: &types.CreateCategoryReq{
				Name:     "Electronics",
				ParentId: 0,
				Sort:     1,
				Icon:     "electronics.png",
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateCategory(
					mock.Anything,
					&product.CreateCategoryRequest{
						Name:     "Electronics",
						ParentId: 0,
						Sort:     1,
						Icon:     "electronics.png",
					},
				).Return(&product.CreateCategoryResponse{
					Id: 1,
				}, nil)
			},
			want: &types.CreateCategoryResp{
				Id: 1,
			},
			wantErr: nil,
		},
		{
			name: "Empty name",
			req: &types.CreateCategoryReq{
				ParentId: 0,
				Sort:     1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateCategory(
					mock.Anything,
					&product.CreateCategoryRequest{
						ParentId: 0,
						Sort:     1,
					},
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			want:    nil,
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Duplicate name",
			req: &types.CreateCategoryReq{
				Name:     "Electronics",
				ParentId: 0,
				Sort:     1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateCategory(
					mock.Anything,
					&product.CreateCategoryRequest{
						Name:     "Electronics",
						ParentId: 0,
						Sort:     1,
					},
				).Return(nil, zeroerr.ErrCategoryDuplicate)
			},
			want:    nil,
			wantErr: zeroerr.ErrCategoryDuplicate,
		},
		{
			name: "Invalid parent",
			req: &types.CreateCategoryReq{
				Name:     "Phones",
				ParentId: 99999,
				Sort:     1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateCategory(
					mock.Anything,
					&product.CreateCategoryRequest{
						Name:     "Phones",
						ParentId: 99999,
						Sort:     1,
					},
				).Return(nil, zeroerr.ErrCategoryNotFound)
			},
			want:    nil,
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

			logic := NewCreateCategoryLogic(context.Background(), svcCtx)
			got, err := logic.CreateCategory(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
