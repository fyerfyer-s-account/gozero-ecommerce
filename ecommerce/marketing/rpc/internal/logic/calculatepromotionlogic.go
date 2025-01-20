package logic

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	// Calculate original amount
	var originalAmount float64
	for _, item := range in.Items {
		originalAmount += item.Price * float64(item.Quantity)
	}

	// Get active promotions with transaction to ensure consistency
	var activePromotions []*model.Promotions
	err := l.svcCtx.PromotionsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		promos, err := l.svcCtx.PromotionsModel.FindActive(ctx)
		if err != nil {
			return err
		}
		activePromotions = promos
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort promotions by type for priority
	sort.Slice(activePromotions, func(i, j int) bool {
		return activePromotions[i].Type < activePromotions[j].Type
	})

	// Calculate discounts
	var (
        bestDiscount   float64
        bestPromotion *model.Promotions
        appliedRule   string
    )

    // Calculate discounts
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
            discountAmount = originalAmount * rule.DiscountRate
            if rule.MaxDiscountAmount > 0 && discountAmount > rule.MaxDiscountAmount {
                discountAmount = rule.MaxDiscountAmount
            }
        }

        // Keep track of best discount
        if discountAmount > bestDiscount {
            bestDiscount = discountAmount
            bestPromotion = promotion
            appliedRule = promotion.Rules
        }
    }

    // Apply only the best promotion
    var appliedPromos []*marketing.PromotionResult
    if bestPromotion != nil {
        appliedPromos = append(appliedPromos, &marketing.PromotionResult{
            PromotionId:    int64(bestPromotion.Id),
            PromotionName:  bestPromotion.Name,
            DiscountAmount: bestDiscount,
        })
    }

    finalAmount := originalAmount - bestDiscount
    if finalAmount < 0 {
        finalAmount = 0
    }

    // Publish calculation event
    event := types.NewMarketingEvent(types.EventTypePromotionCalculated, &types.PromotionCalculateData{
        OrderAmount:    originalAmount,
        DiscountAmount: bestDiscount,
        FinalAmount:    finalAmount,
        AppliedRules:   []string{appliedRule},
        PromotionIds:   []int64{int64(bestPromotion.Id)},
    })

    if err := l.svcCtx.Producer.PublishPromotionEvent(event); err != nil {
        l.Logger.Errorf("Failed to publish promotion calculation event: %v", err)
    }

    return &marketing.CalculatePromotionResponse{
        OriginalAmount:  originalAmount,
        DiscountAmount: bestDiscount,
        FinalAmount:    finalAmount,
        Promotions:     appliedPromos,
    }, nil
}
