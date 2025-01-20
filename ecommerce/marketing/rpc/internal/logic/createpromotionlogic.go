package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreatePromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePromotionLogic {
	return &CreatePromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 促销活动
func (l *CreatePromotionLogic) CreatePromotion(in *marketing.CreatePromotionRequest) (*marketing.CreatePromotionResponse, error) {
    // Validate input
    if err := l.validateInput(in); err != nil {
        return nil, err
    }

    // Validate promotion rules
    var rule types.PromotionRule
    if err := json.Unmarshal([]byte(in.Rules), &rule); err != nil {
        return nil, zeroerr.ErrInvalidPromotionRules
    }

    var promotionId uint64
    err := l.svcCtx.PromotionsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // Create promotion
        result, err := l.svcCtx.PromotionsModel.Insert(ctx, &model.Promotions{
            Name:      in.Name,
            Type:      int64(in.Type),
            Rules:     in.Rules,
            Status:    0, // Not started
            StartTime: sql.NullTime{Time: time.Unix(in.StartTime, 0), Valid: true},
            EndTime:   sql.NullTime{Time: time.Unix(in.EndTime, 0), Valid: true},
        })
        if err != nil {
            return err
        }

        id, err := result.LastInsertId()
        if err != nil {
            return err
        }
        promotionId = uint64(id)

        return nil
    })

    if err != nil {
        return nil, err
    }

    // Publish promotion created event
    event := types.NewMarketingEvent(types.EventTypePromotionCreated, &types.PromotionEventData{
        PromotionID: int64(promotionId),
        Name:        in.Name,
        Type:        in.Type,
        Rules:       in.Rules,
        Status:      0,
        StartTime:   in.StartTime,
        EndTime:     in.EndTime,
    })

    if err := l.svcCtx.Producer.PublishPromotionEvent(event); err != nil {
        l.Logger.Errorf("Failed to publish promotion created event: %v", err)
    }

    return &marketing.CreatePromotionResponse{
        Id: int64(promotionId),
    }, nil
}

func (l *CreatePromotionLogic) validateInput(in *marketing.CreatePromotionRequest) error {
    if in.Name == "" {
        return zeroerr.ErrInvalidMarketingParam
    }
    if in.Type < 1 || in.Type > 3 {
        return zeroerr.ErrInvalidPromotionType
    }
    if in.Rules == "" {
        return zeroerr.ErrInvalidPromotionRules
    }
    if in.StartTime >= in.EndTime {
        return zeroerr.ErrInvalidMarketingParam
    }
    if time.Unix(in.EndTime, 0).Before(time.Now()) {
        return zeroerr.ErrMarketingExpired
    }
    return nil
}
