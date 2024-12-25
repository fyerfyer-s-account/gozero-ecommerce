package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

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
	// todo: add your logic here and delete this line

	return &marketing.GetUserPointsResponse{}, nil
}
