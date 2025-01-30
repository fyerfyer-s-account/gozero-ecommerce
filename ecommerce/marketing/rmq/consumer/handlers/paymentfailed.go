package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PaymentFailedHandler struct {
    logger            *zerolog.Logger
    userCouponsModel  model.UserCouponsModel
    userPointsModel   model.UserPointsModel
    pointsRecordModel model.PointsRecordsModel
}

func NewPaymentFailedHandler(
    userCouponsModel model.UserCouponsModel,
    userPointsModel model.UserPointsModel,
    pointsRecordModel model.PointsRecordsModel,
) *PaymentFailedHandler {
    return &PaymentFailedHandler{
        logger:            zerolog.GetLogger(),
        userCouponsModel:  userCouponsModel,
        userPointsModel:   userPointsModel,
        pointsRecordModel: pointsRecordModel,
    }
}

func (h *PaymentFailedHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.MarketingPaymentFailedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return fmt.Errorf("failed to unmarshal payment failed event: %w", err)
    }

    fields := map[string]interface{}{
        "user_id":    event.UserID,
        "order_no":   event.OrderNo,
        "payment_no": event.PaymentNo,
        "reason":     event.Reason,
    }
    h.logger.Info(ctx, "Processing payment failed event", fields)

    // Handle in transaction
    err := h.userCouponsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Roll back coupon if used
        coupon, err := h.userCouponsModel.FindByOrderNo(ctx, event.OrderNo)
        if err != nil && err != model.ErrNotFound {
            return fmt.Errorf("failed to find coupon by order: %w", err)
        }

        if coupon != nil {
            // Reset coupon to unused state
            if err := h.userCouponsModel.UpdateStatus(ctx, int64(coupon.Id), 0, ""); err != nil {
                return fmt.Errorf("failed to reset coupon status: %w", err)
            }
        }

        // Roll back points if any were awarded
        records, err := h.pointsRecordModel.FindByOrderNo(ctx, event.OrderNo)
        if err != nil && err != model.ErrNotFound {
            return fmt.Errorf("failed to find points record: %w", err)
        }

        if len(records) > 0 {
            for _, record := range records {
                if record.Type == 1 { // Points earned
                    // Create refund record
                    refundRecord := &model.PointsRecords{
                        UserId:  uint64(event.UserID),
                        Points:  -record.Points, // Negative points for refund
                        Type:    2,             // Points used/refunded
                        Source:  "payment_failed",
                        Remark:  sql.NullString{String: fmt.Sprintf("Payment failed refund for order %s", event.OrderNo), Valid: true},
                        OrderNo: sql.NullString{String: event.OrderNo, Valid: true},
                    }

                    if _, err := h.pointsRecordModel.Insert(ctx, refundRecord); err != nil {
                        return fmt.Errorf("failed to create points refund record: %w", err)
                    }

                    // Deduct points from user balance
                    if err := h.userPointsModel.DecrPoints(ctx, event.UserID, record.Points); err != nil {
                        return fmt.Errorf("failed to deduct points: %w", err)
                    }
                }
            }
        }

        return nil
    })

    if err != nil {
        h.logger.Error(ctx, "Failed to process payment failed event", err, fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment failed event", fields)
    return nil
}