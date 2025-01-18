package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListTemplatesLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewListTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTemplatesLogic {
    return &ListTemplatesLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *ListTemplatesLogic) ListTemplates(in *message.ListTemplatesRequest) (*message.ListTemplatesResponse, error) {
    templates, err := l.svcCtx.MessageTemplatesModel.ListByTypeAndStatus(
        l.ctx,
        int64(in.Type),
        int64(in.Status),
        int(in.Page),
        int(in.PageSize),
    )
    if err != nil {
        return nil, zeroerr.ErrMessageNotFound
    }

    total, err := l.svcCtx.MessageTemplatesModel.CountByTypeAndStatus(
        l.ctx,
        int64(in.Type),
        int64(in.Status),
    )
    if err != nil {
        return nil, zeroerr.ErrMessageNotFound
    }

    protoTemplates := make([]*message.MessageTemplate, 0, len(templates))
    for _, tpl := range templates {
        var channels []int32
        if err := json.Unmarshal([]byte(tpl.Channels), &channels); err != nil {
            logx.Errorf("Failed to parse channels for template %d: %v", tpl.Id, err)
            continue
        }

        protoTemplates = append(protoTemplates, &message.MessageTemplate{
            Id:              int64(tpl.Id),
            Code:           tpl.Code,
            Name:           tpl.Name,
            TitleTemplate:  tpl.TitleTemplate,
            ContentTemplate: tpl.ContentTemplate,
            Type:           int32(tpl.Type),
            Channels:       channels,
            Config:         tpl.Config.String,
            Status:         int32(tpl.Status),
            CreatedAt:      tpl.CreatedAt.Unix(),
            UpdatedAt:      tpl.UpdatedAt.Unix(),
        })
    }

    return &message.ListTemplatesResponse{
        Templates: protoTemplates,
        Total:    total,
    }, nil
}