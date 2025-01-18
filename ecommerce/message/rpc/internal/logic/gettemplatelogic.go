package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplateLogic {
	return &GetTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTemplateLogic) GetTemplate(in *message.GetTemplateRequest) (*message.GetTemplateResponse, error) {
	// Get template from database
	template, err := l.svcCtx.MessageTemplatesModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrTemplateNotFound
	}

	// Parse channels from JSON string
	var channels []int32
	err = json.Unmarshal([]byte(template.Channels), &channels)
	if err != nil {
		logx.Errorf("Failed to parse channels: %v", err)
		return nil, zeroerr.ErrInvalidTemplate
	}

	// Convert to proto message
	return &message.GetTemplateResponse{
		Template: &message.MessageTemplate{
			Id:              int64(template.Id),
			Code:            template.Code,
			Name:            template.Name,
			TitleTemplate:   template.TitleTemplate,
			ContentTemplate: template.ContentTemplate,
			Type:            int32(template.Type),
			Channels:        channels,
			Config:          template.Config.String,
			Status:          int32(template.Status),
			CreatedAt:       template.CreatedAt.Unix(),
			UpdatedAt:       template.UpdatedAt.Unix(),
		},
	}, nil
}
