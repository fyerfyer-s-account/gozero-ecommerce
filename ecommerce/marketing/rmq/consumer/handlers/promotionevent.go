package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type PromotionEventHandler struct {
	logger          *zerolog.Logger
	promotionsModel model.PromotionsModel
}

func NewPromotionEventHandler(
	promotionsModel model.PromotionsModel,
) *PromotionEventHandler {
	return &PromotionEventHandler{
		logger:          zerolog.GetLogger(),
		promotionsModel: promotionsModel,
	}
}

func (h *PromotionEventHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	var event types.PromotionEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal promotion event: %w", err)
	}

	fields := map[string]interface{}{
		"type":           event.Type,
		"promotion_id":   event.PromotionID,
		"promotion_type": event.PromotionType,
		"discount":       event.Discount,
		"order_no":       event.OrderNo,
	}
	h.logger.Info(ctx, "Processing promotion event", fields)

	switch event.Type {
	case types.PromotionStarted:
		return h.handlePromotionStarted(ctx, &event)
	case types.PromotionEnded:
		return h.handlePromotionEnded(ctx, &event)
	case types.PromotionApplied:
		return h.handlePromotionApplied(ctx, &event)
	default:
		return fmt.Errorf("unknown promotion event type: %s", event.Type)
	}
}

func (h *PromotionEventHandler) handlePromotionStarted(ctx context.Context, event *types.PromotionEvent) error {
	// Get promotion
	promotion, err := h.promotionsModel.FindOne(ctx, uint64(event.PromotionID))
	if err != nil {
		return fmt.Errorf("failed to get promotion: %w", err)
	}

	// Update promotion status to active
	err = h.promotionsModel.UpdateStatus(ctx, uint64(event.PromotionID), 1) // 1: Active
	if err != nil {
		return fmt.Errorf("failed to update promotion status: %w", err)
	}

	h.logger.Info(ctx, "Promotion started", map[string]interface{}{
		"promotion_id":   event.PromotionID,
		"promotion_name": promotion.Name,
		"start_time":     time.Now(),
	})

	return nil
}

func (h *PromotionEventHandler) handlePromotionEnded(ctx context.Context, event *types.PromotionEvent) error {
	// Get promotion
	promotion, err := h.promotionsModel.FindOne(ctx, uint64(event.PromotionID))
	if err != nil {
		return fmt.Errorf("failed to get promotion: %w", err)
	}

	// Update promotion status to ended
	err = h.promotionsModel.UpdateStatus(ctx, uint64(event.PromotionID), 2) // 2: Ended
	if err != nil {
		return fmt.Errorf("failed to update promotion status: %w", err)
	}

	h.logger.Info(ctx, "Promotion ended", map[string]interface{}{
		"promotion_id":   event.PromotionID,
		"promotion_name": promotion.Name,
		"end_time":       time.Now(),
	})

	return nil
}

func (h *PromotionEventHandler) handlePromotionApplied(ctx context.Context, event *types.PromotionEvent) error {
	// Verify promotion exists and is active
	promotion, err := h.promotionsModel.FindOne(ctx, uint64(event.PromotionID))
	if err != nil {
		return fmt.Errorf("failed to get promotion: %w", err)
	}

	if promotion.Status != 1 { // 1: Active
		return fmt.Errorf("promotion is not active")
	}

	h.logger.Info(ctx, "Promotion applied", map[string]interface{}{
		"promotion_id":   event.PromotionID,
		"promotion_name": promotion.Name,
		"order_no":       event.OrderNo,
		"discount":       event.Discount,
		"applied_time":   time.Now(),
	})

	return nil
}
