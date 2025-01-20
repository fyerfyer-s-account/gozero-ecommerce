package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

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
    if in.Page <= 0 {
        in.Page = 1
    }
    if in.PageSize <= 0 || in.PageSize > 100 {
        in.PageSize = 10
    }

    promotions, err := l.svcCtx.PromotionsModel.FindByStatus(l.ctx, in.Status, in.Page, in.PageSize)
    if err != nil {
        l.Logger.Errorf("Failed to get promotions list: %v", err)
        return nil, zeroerr.ErrPromotionNotFound
    }

    total, err := l.svcCtx.PromotionsModel.Count(l.ctx, in.Status)
    if err != nil {
        l.Logger.Errorf("Failed to get promotions count: %v", err)
        return nil, zeroerr.ErrPromotionNotFound
    }

    var result []*marketing.Promotion
    for _, p := range promotions {
        result = append(result, &marketing.Promotion{
            Id:        int64(p.Id),
            Name:      p.Name,
            Type:      int32(p.Type),
            Rules:     p.Rules,
            Status:    int32(p.Status),
            StartTime: p.StartTime.Time.Unix(),
            EndTime:   p.EndTime.Time.Unix(),
            CreatedAt: p.CreatedAt.Unix(),
            UpdatedAt: p.UpdatedAt.Unix(),
        })
    }

    return &marketing.ListPromotionsResponse{
        Promotions: result,
        Total:     total,
    }, nil
}

