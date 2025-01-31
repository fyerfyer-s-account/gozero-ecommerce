package logic

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/zeromicro/go-zero/core/logx"
)

type UpdateTemplateLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewUpdateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTemplateLogic {
    return &UpdateTemplateLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *UpdateTemplateLogic) UpdateTemplate(in *message.UpdateTemplateRequest) (*message.UpdateTemplateResponse, error) {
    // Find template
    template, err := l.svcCtx.MessageTemplatesModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        return nil, zeroerr.ErrTemplateNotFound
    }

    // Validate status
    if in.Status != 0 && (in.Status < 1 || in.Status > 2) {
        return nil, zeroerr.ErrInvalidParam
    }

    // Convert channels to JSON string
    channelsJSON, err := json.Marshal(in.Channels)
    if err != nil {
        return nil, zeroerr.ErrInvalidParam
    }

    // Update template
    template.Name = in.Name
    template.TitleTemplate = in.TitleTemplate
    template.ContentTemplate = in.ContentTemplate
    template.Channels = string(channelsJSON)
    template.Config = sql.NullString{String: in.Config, Valid: in.Config != ""}
    if in.Status != 0 {
        template.Status = int64(in.Status)
    }

    err = l.svcCtx.MessageTemplatesModel.Update(l.ctx, template)
    if err != nil {
        return nil, zeroerr.ErrMessageUpdateFailed
    }

    // Create template event
    templateEvent := &types.MessageTemplateEvent{
        MessageEvent: types.MessageEvent{
            Type:      types.MessageTemplateUpdated,
            UserID:    0, // System operation
            Timestamp: time.Now(),
        },
        Template: types.MessageTemplate{
            TemplateID:   int64(template.Id),
            TemplateCode: template.Code,
            Content:      template.ContentTemplate,
            Variables:    make(map[string]string), // Initialize empty variables
        },
        Action: "updated",
    }

    // Publish template event
    if err := l.svcCtx.Producer.PublishTemplateEvent(l.ctx, templateEvent); err != nil {
        logx.Errorf("Failed to publish template updated event: %v", err)
        // Don't return error as template is already updated
    }

    return &message.UpdateTemplateResponse{
        Success: true,
    }, nil
}