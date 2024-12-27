package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 地址管理
func (l *AddAddressLogic) AddAddress(in *user.AddAddressRequest) (*user.AddAddressResponse, error) {
	// 1. Check user and validate input
	if err := l.validateUser(in.UserId); err != nil {
		return nil, err
	}
	if err := l.validateAddressInput(in); err != nil {
		return nil, err
	}

	// 2. Check address limit and handle default logic
	_, isDefault, err := l.handleAddressRules(in)
	if err != nil {
		return nil, err
	}

	// 3. Insert new address
	addressId, err := l.insertAddress(in, isDefault)
	if err != nil {
		return nil, err
	}

	return &user.AddAddressResponse{
		AddressId: addressId,
	}, nil
}

func (l *AddAddressLogic) validateUser(userId int64) error {
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(userId))
	if err != nil {
		if err == model.ErrNotFound {
			return zeroerr.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (l *AddAddressLogic) validateAddressInput(in *user.AddAddressRequest) error {
	if len(in.ReceiverName) == 0 || len(in.ReceiverPhone) == 0 ||
		len(in.Province) == 0 || len(in.City) == 0 ||
		len(in.District) == 0 || len(in.DetailAddress) == 0 {
		return zeroerr.ErrInvalidAddress
	}
	return nil
}

// 处理地址逻辑
func (l *AddAddressLogic) handleAddressRules(in *user.AddAddressRequest) ([]*model.UserAddresses, int64, error) {
	addresses, err := l.svcCtx.UserAddressesModel.FindByUserId(l.ctx, uint64(in.UserId))
	if err != nil && err != model.ErrNotFound {
		return nil, 0, err
	}

	if len(addresses) >= l.svcCtx.Config.MaxAddressCount {
		return nil, 0, zeroerr.ErrAddressLimit
	}

	isDefault := int64(0)
	if in.IsDefault || len(addresses) == 0 { // First address is default
		isDefault = 1
	}

	return addresses, isDefault, nil
}

func (l *AddAddressLogic) insertAddress(in *user.AddAddressRequest, isDefault int64) (int64, error) {
	result, err := l.svcCtx.UserAddressesModel.Insert(l.ctx, &model.UserAddresses{
		UserId:        uint64(in.UserId),
		ReceiverName:  in.ReceiverName,
		ReceiverPhone: in.ReceiverPhone,
		Province:      in.Province,
		City:          in.City,
		District:      in.District,
		DetailAddress: in.DetailAddress,
		IsDefault:     isDefault,
	})
	if err != nil {
		logx.Errorf("insert address error: %v", err)
		return 0, zeroerr.ErrInvalidAddress
	}

	addressId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if isDefault == 1 {
		err = l.svcCtx.UserAddressesModel.SetDefault(l.ctx, uint64(in.UserId), uint64(addressId))
		if err != nil {
			logx.Errorf("set default address error: %v", err)
		}
	}

	return addressId, nil
}
