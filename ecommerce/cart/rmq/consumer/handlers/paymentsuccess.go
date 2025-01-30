package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type PaymentSuccessHandler struct {
    logger         *zerolog.Logger
    cartItemsModel model.CartItemsModel
    cartStatsModel model.CartStatisticsModel
}

func NewPaymentSuccessHandler(
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *PaymentSuccessHandler {
    return &PaymentSuccessHandler{
        logger:         zerolog.GetLogger(),
        cartItemsModel: cartItemsModel,
        cartStatsModel: cartStatsModel,
    }
}

func (h *PaymentSuccessHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.CartPaymentSuccessEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal payment success event",
                Err:     err,
            },
        }
    }

    if err := event.Validate(); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "user_id":    event.UserID,
        "order_no":   event.OrderNo,
        "payment_no": event.PaymentNo,
    }
    h.logger.Info(ctx, "Processing payment success event", fields)

    // Delete all cart items for the user
    if err := h.cartItemsModel.DeleteByUserId(ctx, uint64(event.UserID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to clear cart items after payment", fields)
        return err
    }

    // Reset cart statistics
    stats := &model.CartStatistics{
        UserId:           uint64(event.UserID),
        TotalQuantity:    0,
        SelectedQuantity: 0,
        TotalAmount:      0,
        SelectedAmount:   0,
    }

    if err := h.cartStatsModel.Upsert(ctx, stats); err != nil {
        h.logger.WithError(ctx, err, "Failed to reset cart statistics after payment", fields)
        return err
    }

    h.logger.Info(ctx, "Successfully cleared cart after payment", fields)
    return nil
}