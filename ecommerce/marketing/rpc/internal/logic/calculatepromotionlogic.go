package logic

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PromotionRule struct {
	MinAmount         float64 `json:"minAmount"`
	DiscountAmount    float64 `json:"discountAmount"`
	DiscountRate      float64 `json:"discountRate"`
	MaxDiscountAmount float64 `json:"maxDiscountAmount"`
}

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
	if len(in.Items) == 0 {
		return nil, zeroerr.ErrInvalidMarketingParam
	}

	var originalAmount float64
	for _, item := range in.Items {
		originalAmount += item.Price * float64(item.Quantity)
	}

	activePromotions, err := l.svcCtx.PromotionsModel.FindActive(l.ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(activePromotions, func(i, j int) bool {
		return activePromotions[i].Type < activePromotions[j].Type
	})

	var (
		totalDiscount float64
		appliedRules  []string
		promotionIds  []int64
		appliedPromos []*marketing.PromotionResult
	)

	for _, promotion := range activePromotions {
		var rule PromotionRule
		if err := json.Unmarshal([]byte(promotion.Rules), &rule); err != nil {
			l.Logger.Errorf("Invalid promotion rule: %v", err)
			continue
		}

		if originalAmount < rule.MinAmount {
			continue
		}

		var discountAmount float64
		switch promotion.Type {
		case 1: // Fixed amount
			discountAmount = rule.DiscountAmount
		case 2: // Percentage
			discountAmount = originalAmount * (rule.DiscountRate) 
			if rule.MaxDiscountAmount > 0 && discountAmount > rule.MaxDiscountAmount {
				discountAmount = rule.MaxDiscountAmount
			}
		}

		if discountAmount > 0 {
			totalDiscount += discountAmount
			appliedRules = append(appliedRules, promotion.Rules)
			promotionIds = append(promotionIds, int64(promotion.Id))
			appliedPromos = append(appliedPromos, &marketing.PromotionResult{
				PromotionId:    int64(promotion.Id),
				PromotionName:  promotion.Name,
				DiscountAmount: discountAmount,
			})
		}
	}

	finalAmount := originalAmount - totalDiscount
	if finalAmount < 0 {
		finalAmount = 0
	}

	// Publish calculation event
	event := types.NewMarketingEvent(types.EventTypePromotionCalculated, &types.PromotionCalculateData{
		OrderAmount:    originalAmount,
		DiscountAmount: totalDiscount,
		FinalAmount:    finalAmount,
		AppliedRules:   appliedRules,
		PromotionIds:   promotionIds,
	})

	if err := l.svcCtx.Producer.PublishPromotionEvent(event); err != nil {
		l.Logger.Errorf("Failed to publish promotion calculation event: %v", err)
	}

	return &marketing.CalculatePromotionResponse{
		OriginalAmount: originalAmount,
		DiscountAmount: totalDiscount,
		FinalAmount:    finalAmount,
		Promotions:     appliedPromos,
	}, nil
}
