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

func TestCreateSkuLogic_CreateSku(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.CreateSkuReq
		mock    func(mockProduct *ProductService)
		want    *types.CreateSkuResp
		wantErr error
	}{
		{
			name: "Success - Basic SKU",
			req: &types.CreateSkuReq{
				ProductId: 1,
				SkuCode:   "TEST-SKU-001",
				Price:     99.99,
				Stock:     100,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateSku(
					mock.Anything,
					&product.CreateSkuRequest{
						ProductId:  1,
						SkuCode:    "TEST-SKU-001",
						Price:      99.99,
						Stock:      100,
						Attributes: []*product.SkuAttribute{},
					},
				).Return(&product.CreateSkuResponse{Id: 1}, nil)
			},
			want:    &types.CreateSkuResp{Id: 1},
			wantErr: nil,
		},
		{
			name: "Success - With attributes",
			req: &types.CreateSkuReq{
				ProductId: 1,
				SkuCode:   "TEST-SKU-002",
				Price:     199.99,
				Stock:     50,
				Attributes: []types.SkuAttributeReq{
					{Key: "Color", Value: "Red"},
					{Key: "Size", Value: "L"},
				},
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateSku(
					mock.Anything,
					&product.CreateSkuRequest{
						ProductId: 1,
						SkuCode:   "TEST-SKU-002",
						Price:     199.99,
						Stock:     50,
						Attributes: []*product.SkuAttribute{
							{Key: "Color", Value: "Red"},
							{Key: "Size", Value: "L"},
						},
					},
				).Return(&product.CreateSkuResponse{Id: 2}, nil)
			},
			want:    &types.CreateSkuResp{Id: 2},
			wantErr: nil,
		},
		{
			name: "Product not found",
			req: &types.CreateSkuReq{
				ProductId: 999,
				SkuCode:   "TEST-SKU-003",
				Price:     99.99,
				Stock:     100,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateSku(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrProductNotFound)
			},
			want:    nil,
			wantErr: zeroerr.ErrProductNotFound,
		},
		{
			name: "Duplicate SKU code",
			req: &types.CreateSkuReq{
				ProductId: 1,
				SkuCode:   "EXISTING-SKU",
				Price:     99.99,
				Stock:     100,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateSku(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrSkuDuplicate)
			},
			want:    nil,
			wantErr: zeroerr.ErrSkuDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewCreateSkuLogic(context.Background(), svcCtx)
			got, err := logic.CreateSku(tt.req)

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
