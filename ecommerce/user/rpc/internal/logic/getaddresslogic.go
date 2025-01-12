package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressLogic {
	return &GetAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAddressLogic) GetAddress(in *user.GetAddressRequest) (*user.GetAddressResponse, error) {
	// todo: add your logic here and delete this line
	addr, err := l.svcCtx.UserAddressesModel.FindOne(l.ctx, uint64(in.AddressId)) 
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrUserNotFound
		}
		return nil, err
	}

	return &user.GetAddressResponse{
		Address: &user.Address{
			Id: in.AddressId,
			UserId: int64(addr.UserId),
			ReceiverName: addr.ReceiverName,
			ReceiverPhone: addr.ReceiverPhone,
			Province: addr.Province,
			City: addr.City,
			District: addr.District,
			DetailAddress: addr.DetailAddress,
			IsDefault: addr.IsDefault == 1,
			CreatedAt: addr.CreatedAt.Unix(),
		},
	}, nil
}
