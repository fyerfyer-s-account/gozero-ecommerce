package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAddressesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAddressesLogic {
	return &GetUserAddressesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserAddressesLogic) GetUserAddresses(in *user.GetUserAddressesRequest) (*user.GetUserAddressesResponse, error) {
	// todo: add your logic here and delete this line
	// 1. Check if user exists
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrUserNotFound
		}
		return nil, err
	}

	// 2. Get user addresses
	addresses, err := l.svcCtx.UserAddressesModel.FindByUserId(l.ctx, uint64(in.UserId))
	if err != nil {
		logx.Errorf("get user addresses error: %v", err)
		return nil, zeroerr.ErrAddressNotFound
	}

	// 3. Convert model data to proto format
	resp := make([]*user.Address, 0, len(addresses))
	for _, addr := range addresses {
		resp = append(resp, &user.Address{
			Id:            int64(addr.Id),
			UserId:        int64(addr.UserId),
			ReceiverName:  addr.ReceiverName,
			ReceiverPhone: addr.ReceiverPhone,
			Province:      addr.Province,
			City:          addr.City,
			District:      addr.District,
			DetailAddress: addr.DetailAddress,
			IsDefault:     addr.IsDefault == 1,
			CreatedAt:     addr.CreatedAt.Unix(),
			UpdatedAt:     addr.UpdatedAt.Unix(),
		})
	}

	return &user.GetUserAddressesResponse{
		Addresses: resp,
	}, nil
}
