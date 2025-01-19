package handlers

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type PromotionHandler struct {
    promotionsModel model.PromotionsModel
}

func NewPromotionHandler(promotionsModel model.PromotionsModel) *PromotionHandler {
    return &PromotionHandler{
        promotionsModel: promotionsModel,
    }
}

func (h *PromotionHandler) HandlePromotionStatus(event *types.MarketingEvent) error {
    data, ok := event.Data.(*types.PromotionEventData)
    if !ok {
        return zeroerr.ErrInvalidEventData
    }

    promotion, err := h.promotionsModel.FindOne(context.Background(), uint64(data.PromotionID))
    if err != nil {
        return err
    }

    // Update promotion status
    promotion.Status = int64(data.Status)
    if err := h.promotionsModel.Update(context.Background(), promotion); err != nil {
        return err
    }

    return nil
}

func (h *PromotionHandler) HandlePromotionCalculated(event *types.MarketingEvent) error {
    data, ok := event.Data.(*types.PromotionCalculateData)
    if !ok {
        return zeroerr.ErrInvalidEventData
    }

    // Log promotion calculation
    logx.Infof("Promotion calculated: OrderAmount=%.2f, DiscountAmount=%.2f, FinalAmount=%.2f",
        data.OrderAmount, data.DiscountAmount, data.FinalAmount)

    return nil
}