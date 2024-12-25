package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPromotionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPromotionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPromotionsLogic {
	return &ListPromotionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPromotionsLogic) ListPromotions(in *marketing.ListPromotionsRequest) (*marketing.ListPromotionsResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.ListPromotionsResponse{}, nil
}
