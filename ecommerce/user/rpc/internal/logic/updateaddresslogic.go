package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAddressLogic) UpdateAddress(in *user.UpdateAddressRequest) (*user.UpdateAddressResponse, error) {
	// 1. Validate input
	if err := l.validateAddressInput(in); err != nil {
		return nil, err
	}

	// 2. Get and validate address
	address, err := l.getAndValidateAddress(in.AddressId)
	if err != nil {
		return nil, err
	}

	// 3. Update address with transaction
	err = l.svcCtx.UserAddressesModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Update address info
		err := l.svcCtx.UserAddressesModel.WithSession(session).Update(ctx, &model.UserAddresses{
			Id:            uint64(in.AddressId),
			ReceiverName:  in.ReceiverName,
			ReceiverPhone: in.ReceiverPhone,
			Province:      in.Province,
			City:          in.City,
			District:      in.District,
			DetailAddress: in.DetailAddress,
			IsDefault:     boolToInt64(in.IsDefault),
		})
		if err != nil {
			return err
		}

		// Handle default address logic if needed
		if in.IsDefault && address.IsDefault == 0 {
			err = l.svcCtx.UserAddressesModel.WithSession(session).SetDefault(ctx, address.UserId, uint64(in.AddressId))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logx.Errorf("update address error: %v", err)
		return nil, err
	}

	return &user.UpdateAddressResponse{
		Success: true,
	}, nil
}

func (l *UpdateAddressLogic) validateAddressInput(in *user.UpdateAddressRequest) error {
	if len(in.ReceiverName) == 0 || len(in.ReceiverPhone) == 0 ||
		len(in.Province) == 0 || len(in.City) == 0 ||
		len(in.District) == 0 || len(in.DetailAddress) == 0 {
		return zeroerr.ErrInvalidAddress
	}
	return nil
}

func (l *UpdateAddressLogic) getAndValidateAddress(addressId int64) (*model.UserAddresses, error) {
	address, err := l.svcCtx.UserAddressesModel.FindOne(l.ctx, uint64(addressId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrAddressNotFound
		}
		return nil, err
	}
	return address, nil
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
