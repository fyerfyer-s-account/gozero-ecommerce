package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UsePointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUsePointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UsePointsLogic {
	return &UsePointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UsePointsLogic) UsePoints(in *marketing.UsePointsRequest) (*marketing.UsePointsResponse, error) {
    // Validate input
    if in.UserId <= 0 || in.Points <= 0 || in.Usage == "" || in.OrderNo == "" {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    var currentPoints int64
    err := l.svcCtx.UserPointsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // Lock user points record
        if err := l.svcCtx.UserPointsModel.Lock(ctx, session, in.UserId); err != nil {
            return err
        }

        // Check points balance
        balance, err := l.svcCtx.UserPointsModel.GetBalance(ctx, in.UserId)
        if err != nil {
            return err
        }
        if balance < in.Points {
            return zeroerr.ErrInsufficientPoints
        }

        // Deduct points
        if err := l.svcCtx.UserPointsModel.DecrPoints(ctx, in.UserId, in.Points); err != nil {
            return err
        }

        // Create points record
        _, err = l.svcCtx.PointsRecordsModel.Insert(ctx, &model.PointsRecords{
            UserId:    uint64(in.UserId),
            Points:    -in.Points, // Negative for points used
            Type:      2,         // Use points
            Source:    in.Usage,
            Remark:    sql.NullString{String: "Points usage", Valid: true},
            OrderNo:   sql.NullString{String: in.OrderNo, Valid: true},
            CreatedAt: time.Now(),
        })
        if err != nil {
            return err
        }

        // Get updated points balance
        currentPoints, err = l.svcCtx.UserPointsModel.GetBalance(ctx, in.UserId)
        return err
    })

    if err != nil {
        l.Logger.Errorf("Failed to use points: %v", err)
        return nil, err
    }

    // Publish points used event
    pointsEvent := &types.PointsEvent{
        MarketingEvent: types.MarketingEvent{
            Type:      types.PointsUsed,
            UserID:    in.UserId,
            Timestamp: time.Now(),
        },
        Points:  in.Points,
        Balance: currentPoints,
        Source:  in.Usage,
        OrderNo: in.OrderNo,
        Reason:  "Points usage",
    }

    if err := l.svcCtx.Producer.PublishPointsEvent(l.ctx, pointsEvent); err != nil {
        l.Logger.Errorf("Failed to publish points used event: %v", err)
        // Don't return error as the main transaction succeeded
    }

    return &marketing.UsePointsResponse{
        Success:       true,
        CurrentPoints: currentPoints,
    }, nil
}