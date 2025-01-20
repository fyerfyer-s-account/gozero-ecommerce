package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPointsLogic {
	return &GetUserPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 积分系统
func (l *GetUserPointsLogic) GetUserPoints(in *marketing.GetUserPointsRequest) (*marketing.GetUserPointsResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    points, err := l.svcCtx.UserPointsModel.GetBalance(l.ctx, in.UserId)
    if err != nil {
        if err == zeroerr.ErrNotFound {
            return &marketing.GetUserPointsResponse{
                Points: 0,
            }, nil
        }
        l.Logger.Errorf("Failed to get user points: %v", err)
        return nil, err
    }

    return &marketing.GetUserPointsResponse{
        Points: points,
    }, nil
}
