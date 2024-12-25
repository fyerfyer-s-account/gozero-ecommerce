package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type PointsHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPointsHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointsHistoryLogic {
	return &PointsHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PointsHistoryLogic) PointsHistory(in *marketing.PointsHistoryRequest) (*marketing.PointsHistoryResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.PointsHistoryResponse{}, nil
}
