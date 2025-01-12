package order

import (
    "context"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cartclient"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreateOrderLogic_CreateOrder(t *testing.T) {
	mockOrder := NewOrderService(t)
	mockCart := NewCart(t)
	mockUser := NewUser(t)
	mockProduct := NewProductService(t)

	svcCtx := &svc.ServiceContext{
		OrderRpc:   mockOrder,
		CartRpc:    mockCart,
		UserRpc:    mockUser,
		ProductRpc: mockProduct,
	}

	tests := []struct {
		name    string
		ctx     context.Context
		req     *types.CreateOrderReq
		mock    func()
		want    *types.Order
		wantErr error
	}{
		{
			name: "create order successfully",
			ctx:  context.WithValue(context.Background(), "userId", int64(1)),
			req: &types.CreateOrderReq{
				AddressId: 1,
			},
			mock: func() {
				mockUser.On("GetAddress",
					mock.Anything,
					&user.GetAddressRequest{AddressId: 1},
				).Return(&user.GetAddressResponse{
					Address: &user.Address{
						DetailAddress:  "Test Address",
						ReceiverName:   "Test User",
						ReceiverPhone:  "1234567890",
					},
				}, nil)

				mockCart.On("GetSelectedItems",
					mock.Anything,
					&cartclient.GetSelectedItemsRequest{UserId: 1},
				).Return(&cartclient.GetSelectedItemsResponse{
					Items: []*cartclient.CartItem{{
						ProductId: 1,
						SkuId:    1,
						Quantity: 2,
					}},
				}, nil)

				mockProduct.On("GetSku",
					mock.Anything,
					&product.GetSkuRequest{Id: 1},
				).Return(&product.GetSkuResponse{}, nil)

				mockOrder.On("CreateOrder",
					mock.Anything,
					&order.CreateOrderRequest{
						UserId:   1,
						Address:  "Test Address",
						Receiver: "Test User",
						Phone:    "1234567890",
						Items: []*order.OrderItemRequest{{
							ProductId: 1,
							SkuId:    1,
							Quantity: 2,
						}},
					},
				).Return(&order.CreateOrderResponse{
					OrderNo: "ORDER123",
				}, nil)

				mockOrder.On("GetOrder",
					mock.Anything,
					&order.GetOrderRequest{OrderNo: "ORDER123"},
				).Return(&order.GetOrderResponse{
					Order: &order.Order{
						Id:      1,
						OrderNo: "ORDER123",
						UserId:  1,
						Status:  1,
					},
				}, nil)
			},
			want: &types.Order{
				Id:      1,
				OrderNo: "ORDER123",
				UserId:  1,
				Status:  1,
			},
			wantErr: nil,
		},
		{
			name: "empty cart items",
			ctx:  context.WithValue(context.Background(), "userId", int64(1)),
			req: &types.CreateOrderReq{
				AddressId: 1,
				Note: "hello",
			},
			mock: func() {
				mockCart.On("GetSelectedItems",
					mock.Anything,
					&cartclient.GetSelectedItemsRequest{UserId: 1},
				).Return(&cartclient.GetSelectedItemsResponse{
					Items: []*cartclient.CartItem{},
				}, nil)
				
				// No need to mock other calls as we should return early
			},
			want:    nil,
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			l := NewCreateOrderLogic(tt.ctx, svcCtx)
			got, err := l.CreateOrder(tt.req)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.OrderNo, got.OrderNo)
				assert.Equal(t, tt.want.UserId, got.UserId)
				assert.Equal(t, tt.want.Status, got.Status)
			}
		})
	}
}
