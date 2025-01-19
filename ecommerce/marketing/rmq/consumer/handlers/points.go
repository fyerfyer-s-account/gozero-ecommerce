package handlers

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PointsHandler struct {
    pointsModel      model.UserPointsModel
    pointsRecordModel model.PointsRecordsModel
}

func NewPointsHandler(pointsModel model.UserPointsModel, recordModel model.PointsRecordsModel) *PointsHandler {
    return &PointsHandler{
        pointsModel:      pointsModel,
        pointsRecordModel: recordModel,
    }
}

func (h *PointsHandler) HandlePointsTransaction(event *types.MarketingEvent) error {
    data, ok := event.Data.(*types.PointsEventData)
    if !ok {
        return zeroerr.ErrInvalidEventData
    }

    return h.pointsModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        // Lock user points
        if err := h.pointsModel.Lock(ctx, session, data.UserID); err != nil {
            return err
        }

        // Add or deduct points
        if data.Type == 1 { // add points
            if err := h.pointsModel.IncrPoints(ctx, data.UserID, data.Points); err != nil {
                return err
            }
        } else { // use points
            if err := h.pointsModel.DecrPoints(ctx, data.UserID, data.Points); err != nil {
                return err
            }
        }

        // Create points record
        _, err := h.pointsRecordModel.Insert(ctx, &model.PointsRecords{
            UserId:  uint64(data.UserID),
            Points:  data.Points,
            Type:    int64(data.Type),
            Source:  data.Source,
            Remark:  sql.NullString{String: data.Remark, Valid: true},
            OrderNo: sql.NullString{String: data.OrderNo, Valid: true},
        })
        return err
    })
}