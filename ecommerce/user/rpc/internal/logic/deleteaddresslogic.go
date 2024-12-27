package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAddressLogic) DeleteAddress(in *user.DeleteAddressRequest) (*user.DeleteAddressResponse, error) {
	if err := l.validateUser(in.UserId); err != nil {
		return nil, err
	}

	if err := l.getAndValidateAddress(in.UserId, in.AddressId); err != nil {
		return nil, err
	}

	if err := l.deleteAddress(in.AddressId); err != nil {
		return nil, err
	}

	return &user.DeleteAddressResponse{
		Success: true,
	}, nil
}

func (l *DeleteAddressLogic) validateUser(userId int64) error {
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(userId))
	if err != nil {
		if err == model.ErrNotFound {
			return zeroerr.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (l *DeleteAddressLogic) getAndValidateAddress(userId, addressId int64) error {
	address, err := l.svcCtx.UserAddressesModel.FindOne(l.ctx, uint64(addressId))
	if err != nil {
		if err == model.ErrNotFound {
			return zeroerr.ErrAddressNotFound
		}
		return err
	}

	if address.UserId != uint64(userId) {
		return zeroerr.ErrInvalidAddress
	}

	if address.IsDefault == 1 {
		return zeroerr.ErrDefaultAddressNotDeletable
	}

	return nil
}

func (l *DeleteAddressLogic) deleteAddress(addressId int64) error {
	err := l.svcCtx.UserAddressesModel.Delete(l.ctx, uint64(addressId))
	if err != nil {
		logx.Errorf("delete address error: %v", err)
		return err
	}
	return nil
}
