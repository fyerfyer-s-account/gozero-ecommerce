package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePromotionLogic {
	return &CreatePromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 促销活动
func (l *CreatePromotionLogic) CreatePromotion(in *marketing.CreatePromotionRequest) (*marketing.CreatePromotionResponse, error) {
	// todo: add your logic here and delete this line

	return &marketing.CreatePromotionResponse{}, nil
}
