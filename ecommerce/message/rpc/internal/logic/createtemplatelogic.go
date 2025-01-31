package logic

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/zeromicro/go-zero/core/logx"
)

type CreateTemplateLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewCreateTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTemplateLogic {
    return &CreateTemplateLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *CreateTemplateLogic) CreateTemplate(in *message.CreateTemplateRequest) (*message.CreateTemplateResponse, error) {
    // Validate input
    if in.Code == "" || in.Name == "" || in.TitleTemplate == "" || in.ContentTemplate == "" {
        return nil, zeroerr.ErrInvalidTemplate
    }

    // Check if template code exists
    _, err := l.svcCtx.MessageTemplatesModel.FindByCode(l.ctx, in.Code)
    if err == nil {
        return nil, zeroerr.ErrDuplicateTemplate
    } else if err != model.ErrNotFound {
        logx.Errorf("Failed to check template code: %v", err)
        return nil, zeroerr.ErrTemplateCreateFailed
    }

    // Convert channels to JSON string
    channelsJson, err := json.Marshal(in.Channels)
    if err != nil {
        return nil, err
    }

    // Create template record
    template := &model.MessageTemplates{
        Code:            in.Code,
        Name:            in.Name,
        TitleTemplate:   in.TitleTemplate,
        ContentTemplate: in.ContentTemplate,
        Type:            int64(in.Type),
        Channels:        string(channelsJson),
        Status:          1, // Active by default
    }

    if in.Config != "" {
        template.Config = sql.NullString{
            String: in.Config,
            Valid:  true,
        }
    }

    result, err := l.svcCtx.MessageTemplatesModel.Insert(l.ctx, template)
    if err != nil {
        logx.Errorf("Failed to create template: %v", err)
        return nil, zeroerr.ErrTemplateCreateFailed
    }

    templateId, err := result.LastInsertId()
    if err != nil {
        logx.Errorf("Failed to get template ID: %v", err)
        return nil, err
    }

    // Create and publish template event
    templateEvent := &types.MessageTemplateEvent{
        MessageEvent: types.MessageEvent{
            Type:      types.MessageTemplateCreated,
            UserID:    0, // System operation
            Timestamp: time.Now(),
        },
        Template: types.MessageTemplate{
            TemplateID:   templateId,
            TemplateCode: in.Code,
            Content:      in.ContentTemplate,
            Variables:    make(map[string]string), // Initialize empty variables
        },
        Action: "created",
    }

    if err := l.svcCtx.Producer.PublishTemplateEvent(l.ctx, templateEvent); err != nil {
        logx.Errorf("Failed to publish template created event: %v", err)
        // Don't return error as template is already created
    }

    return &message.CreateTemplateResponse{
        TemplateId: templateId,
    }, nil
}