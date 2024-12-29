package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RechargeWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRechargeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeWalletLogic {
	return &RechargeWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RechargeWalletLogic) RechargeWallet(req *types.RechargeReq) error {
	// todo: add your logic here and delete this line
	// Get userId from JWT context
	userId := l.ctx.Value("userId").(int64)

	// Call RPC
	_, err := l.svcCtx.UserRpc.RechargeWallet(l.ctx, &user.RechargeWalletRequest{
		UserId:  userId,
		Amount:  req.Amount,
		Channel: getPaymentType(req.PaymentType),
	})

	if err != nil {
		logx.Errorf("recharge wallet error: %v", err)
		return zeroerr.ErrRechargeWalletFailed
	}

	return nil
}

func getPaymentType(channel int32) string {
	switch channel {
	case 1:
		return "alipay"
	case 2:
		return "wechat"
	default:
		return "unknown"
	}
}
