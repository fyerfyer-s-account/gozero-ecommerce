package logic

import (
	"context"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AddPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPointsLogic {
	return &AddPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddPointsLogic) AddPoints(in *marketing.AddPointsRequest) (*marketing.AddPointsResponse, error) {
    if in.Points <= 0 {
        return nil, zeroerr.ErrInvalidPointsAmount
    }

    if in.Points > l.svcCtx.Config.PointsLimits.MaxPoints {
        return nil, zeroerr.ErrExceedPointsLimit
    }

    var currentPoints int64
    err := l.svcCtx.UserPointsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // Always try to initialize first for new users
        balance, err := l.svcCtx.UserPointsModel.GetBalance(ctx, in.UserId)
        if err == zeroerr.ErrNotFound {
            if err := l.svcCtx.UserPointsModel.InitUserPoints(ctx, in.UserId); err != nil {
                return err
            }
            balance = 0
        } else if err != nil {
            return err
        }

        // Add points
        if err := l.svcCtx.UserPointsModel.IncrPoints(ctx, in.UserId, in.Points); err != nil {
            return err
        }

        currentPoints = balance + in.Points
        return nil
    })

    if err != nil {
        return nil, err
    }

    // Publish event
    event := types.NewMarketingEvent(types.EventTypePointsAdded, &types.PointsEventData{
        UserID:  in.UserId,
        Points:  in.Points,
        Type:    1,
        Source:  in.Source,
        Remark:  in.Remark,
    })

    if err := l.svcCtx.Producer.PublishPointsEvent(event); err != nil {
        logx.Errorf("Failed to publish points event: %v", err)
    }

    return &marketing.AddPointsResponse{
        Success:       true,
        CurrentPoints: currentPoints,
    }, nil
}
