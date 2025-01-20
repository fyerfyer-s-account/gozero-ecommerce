package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PointsHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPointsHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointsHistoryLogic {
	return &PointsHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PointsHistoryLogic) PointsHistory(in *marketing.PointsHistoryRequest) (*marketing.PointsHistoryResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidMarketingParam
    }

    if in.Page <= 0 {
        in.Page = 1
    }
    if in.PageSize <= 0 || in.PageSize > 100 {
        in.PageSize = 10
    }

    records, err := l.svcCtx.PointsRecordsModel.FindByUserId(l.ctx, in.UserId, in.Page, in.PageSize)
    if err != nil {
        l.Logger.Errorf("Failed to get points history: %v", err)
        return nil, zeroerr.ErrPointsNotFound
    }

    total, err := l.svcCtx.PointsRecordsModel.CountByUserId(l.ctx, in.UserId)
    if err != nil {
        l.Logger.Errorf("Failed to get points history count: %v", err)
        return nil, zeroerr.ErrPointsNotFound
    }

    var result []*marketing.PointsRecord
    for _, r := range records {
        result = append(result, &marketing.PointsRecord{
            Id:        int64(r.Id),
            UserId:    int64(r.UserId),
            Points:    r.Points,
            Type:      int32(r.Type),
            Source:    r.Source,
            Remark:    r.Remark.String,
            CreatedAt: r.CreatedAt.Unix(),
        })
    }

    return &marketing.PointsHistoryResponse{
        Records: result,
        Total:   total,
    }, nil
}
