package product

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProductSkusLogic_GetProductSkus(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.GetProductSkusReq
		mock    func(mockProduct *ProductService)
		want    []types.Sku
		wantErr error
	}{
		{
			name: "successful get sku",
			req: &types.GetProductSkusReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().GetSku(
					mock.Anything,
					&product.GetSkuRequest{
						Id: 1,
					},
				).Return(&product.GetSkuResponse{
					Sku: &product.Sku{
						Id:        1,
						ProductId: 100,
						SkuCode:   "SKU001",
						Price:     99.99,
						Stock:     100,
						Attributes: []*product.SkuAttribute{
							{Key: "Color", Value: "Red"},
							{Key: "Size", Value: "L"},
						},
					},
				}, nil)
			},
			want: []types.Sku{
				{
					Id:        1,
					ProductId: 100,
					Name:      "SKU001",
					Code:      "SKU001",
					Price:     99.99,
					Stock:     100,
					Attributes: map[string]string{
						"Color": "Red",
						"Size":  "L",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "rpc error",
			req: &types.GetProductSkusReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().GetSku(
					mock.Anything,
					&product.GetSkuRequest{
						Id: 1,
					},
				).Return(nil, errors.New("sku not found"))
			},
			want:    nil,
			wantErr: errors.New("sku not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewGetProductSkusLogic(context.Background(), svcCtx)
			got, err := logic.GetProductSkus(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
