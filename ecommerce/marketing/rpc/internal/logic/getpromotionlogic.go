package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromotionLogic {
	return &GetPromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPromotionLogic) GetPromotion(in *marketing.GetPromotionRequest) (*marketing.GetPromotionResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.GetPromotionResponse{}, nil
}
