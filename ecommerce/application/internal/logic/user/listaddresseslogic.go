package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAddressesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressesLogic {
	return &ListAddressesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAddressesLogic) ListAddresses() (resp []types.Address, err error) {
	// todo: add your logic here and delete this line
	// Get userId from JWT context
	userId := l.ctx.Value("userId").(int64)

	// Call RPC
	addresses, err := l.svcCtx.UserRpc.GetUserAddresses(l.ctx, &user.GetUserAddressesRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("list addresses error: %v", err)
		return nil, err
	}

	result := make([]types.Address, 0, len(addresses.Addresses))
	for _, addr := range addresses.Addresses {
		result = append(result, types.Address{
			Id:            addr.Id,
			ReceiverName:  addr.ReceiverName,
			ReceiverPhone: addr.ReceiverPhone,
			Province:      addr.Province,
			City:          addr.City,
			District:      addr.District,
			DetailAddress: addr.DetailAddress,
			IsDefault:     addr.IsDefault,
		})
	}

	return result, nil
}
