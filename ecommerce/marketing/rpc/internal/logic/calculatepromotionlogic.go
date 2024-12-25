package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculatePromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculatePromotionLogic {
	return &CalculatePromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CalculatePromotionLogic) CalculatePromotion(in *marketing.CalculatePromotionRequest) (*marketing.CalculatePromotionResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.CalculatePromotionResponse{}, nil
}
