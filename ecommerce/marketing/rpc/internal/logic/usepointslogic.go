package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type UsePointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUsePointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UsePointsLogic {
	return &UsePointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UsePointsLogic) UsePoints(in *marketing.UsePointsRequest) (*marketing.UsePointsResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.UsePointsResponse{}, nil
}
