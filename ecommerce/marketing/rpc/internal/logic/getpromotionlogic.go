package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

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
    if in.Id <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    promotion, err := l.svcCtx.PromotionsModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        l.Logger.Errorf("Failed to get promotion: %v", err)
        return nil, zeroerr.ErrPromotionNotFound
    }

    return &marketing.GetPromotionResponse{
        Promotion: &marketing.Promotion{
            Id:        int64(promotion.Id),
            Name:      promotion.Name,
            Type:      int32(promotion.Type),
            Rules:     promotion.Rules,
            Status:    int32(promotion.Status),
            StartTime: promotion.StartTime.Time.Unix(),
            EndTime:   promotion.EndTime.Time.Unix(),
            CreatedAt: promotion.CreatedAt.Unix(),
            UpdatedAt: promotion.UpdatedAt.Unix(),
        },
    }, nil
}
