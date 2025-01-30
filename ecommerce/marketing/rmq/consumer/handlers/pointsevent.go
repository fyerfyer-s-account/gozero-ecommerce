package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PointsEventHandler struct {
    logger            *zerolog.Logger
    userPointsModel   model.UserPointsModel
    pointsRecordModel model.PointsRecordsModel
}

func NewPointsEventHandler(
    userPointsModel model.UserPointsModel,
    pointsRecordModel model.PointsRecordsModel,
) *PointsEventHandler {
    return &PointsEventHandler{
        logger:            zerolog.GetLogger(),
        userPointsModel:   userPointsModel,
        pointsRecordModel: pointsRecordModel,
    }
}

func (h *PointsEventHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.PointsEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return fmt.Errorf("failed to unmarshal points event: %w", err)
    }

    fields := map[string]interface{}{
        "type":     event.Type,
        "user_id":  event.UserID,
        "points":   event.Points,
        "source":   event.Source,
        "order_no": event.OrderNo,
    }
    h.logger.Info(ctx, "Processing points event", fields)

    switch event.Type {
    case types.PointsEarned:
        return h.handlePointsEarned(ctx, &event)
    case types.PointsUsed:
        return h.handlePointsUsed(ctx, &event)
    case types.PointsExpired:
        return h.handlePointsExpired(ctx, &event)
    case types.PointsRefunded:
        return h.handlePointsRefunded(ctx, &event)
    default:
        return fmt.Errorf("unknown points event type: %s", event.Type)
    }
}

func (h *PointsEventHandler) handlePointsEarned(ctx context.Context, event *types.PointsEvent) error {
    // Begin transaction
    err := h.userPointsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Initialize user points if not exists
        _, err := h.userPointsModel.FindOne(ctx, uint64(event.UserID))
        if err == model.ErrNotFound {
            if err := h.userPointsModel.InitUserPoints(ctx, event.UserID); err != nil {
                return fmt.Errorf("failed to init user points: %w", err)
            }
        } else if err != nil {
            return fmt.Errorf("failed to find user points: %w", err)
        }

        // Create points record
        record := &model.PointsRecords{
            UserId:    uint64(event.UserID),
            Points:    event.Points,
            Type:      1, // Points earned
            Source:    event.Source,
            OrderNo:   sql.NullString{String: event.OrderNo, Valid: event.OrderNo != ""},
            Remark:    sql.NullString{String: event.Reason, Valid: event.Reason != ""},
            CreatedAt: time.Now(),
        }

        if _, err := h.pointsRecordModel.Insert(ctx, record); err != nil {
            return fmt.Errorf("failed to create points record: %w", err)
        }

        // Update user points
        if err := h.userPointsModel.IncrPoints(ctx, event.UserID, event.Points); err != nil {
            return fmt.Errorf("failed to increment points: %w", err)
        }

        return nil
    })

    return err
}

func (h *PointsEventHandler) handlePointsUsed(ctx context.Context, event *types.PointsEvent) error {
    err := h.userPointsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Verify points balance
        balance, err := h.userPointsModel.GetBalance(ctx, event.UserID)
        if err != nil {
            return fmt.Errorf("failed to get points balance: %w", err)
        }

        if balance < event.Points {
            return fmt.Errorf("insufficient points balance: %d < %d", balance, event.Points)
        }

        // Create points record
        record := &model.PointsRecords{
            UserId:    uint64(event.UserID),
            Points:    -event.Points, // Negative for points used
            Type:      2, // Points used
            Source:    event.Source,
            OrderNo:   sql.NullString{String: event.OrderNo, Valid: event.OrderNo != ""},
            Remark:    sql.NullString{String: event.Reason, Valid: event.Reason != ""},
            CreatedAt: time.Now(),
        }

        if _, err := h.pointsRecordModel.Insert(ctx, record); err != nil {
            return fmt.Errorf("failed to create points record: %w", err)
        }

        // Update user points
        if err := h.userPointsModel.DecrPoints(ctx, event.UserID, event.Points); err != nil {
            return fmt.Errorf("failed to decrement points: %w", err)
        }

        return nil
    })

    return err
}

func (h *PointsEventHandler) handlePointsExpired(ctx context.Context, event *types.PointsEvent) error {
    err := h.userPointsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Create expiry record
        record := &model.PointsRecords{
            UserId:    uint64(event.UserID),
            Points:    -event.Points, // Negative for expired points
            Type:      3, // Points expired
            Source:    event.Source,
            Remark:    sql.NullString{String: "Points expired", Valid: true},
            CreatedAt: time.Now(),
        }

        if _, err := h.pointsRecordModel.Insert(ctx, record); err != nil {
            return fmt.Errorf("failed to create points record: %w", err)
        }

        // Update user points
        if err := h.userPointsModel.DecrPoints(ctx, event.UserID, event.Points); err != nil {
            return fmt.Errorf("failed to decrement expired points: %w", err)
        }

        return nil
    })

    return err
}

func (h *PointsEventHandler) handlePointsRefunded(ctx context.Context, event *types.PointsEvent) error {
    err := h.userPointsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Create refund record
        record := &model.PointsRecords{
            UserId:    uint64(event.UserID),
            Points:    event.Points,
            Type:      4, // Points refunded
            Source:    event.Source,
            OrderNo:   sql.NullString{String: event.OrderNo, Valid: event.OrderNo != ""},
            Remark:    sql.NullString{String: event.Reason, Valid: event.Reason != ""},
            CreatedAt: time.Now(),
        }

        if _, err := h.pointsRecordModel.Insert(ctx, record); err != nil {
            return fmt.Errorf("failed to create points record: %w", err)
        }

        // Update user points
        if err := h.userPointsModel.IncrPoints(ctx, event.UserID, event.Points); err != nil {
            return fmt.Errorf("failed to increment refunded points: %w", err)
        }

        return nil
    })

    return err
}