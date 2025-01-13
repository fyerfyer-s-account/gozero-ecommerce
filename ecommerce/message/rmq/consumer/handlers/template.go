package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
)

type TemplateHandler struct {
	templates model.MessageTemplatesModel
}

func NewTemplateHandler(templates model.MessageTemplatesModel) *TemplateHandler {
	return &TemplateHandler{
		templates: templates,
	}
}

func (h *TemplateHandler) Handle(event *types.MessageEvent) error {
	switch event.Type {
	case types.EventTypeTemplateCreated:
		return h.handleTemplateCreated(event)
	case types.EventTypeTemplateUpdated:
		return h.handleTemplateUpdated(event)
	default:
		return nil
	}
}

func (h *TemplateHandler) handleTemplateCreated(event *types.MessageEvent) error {
	data, ok := event.Data.(*types.TemplateData)
	if !ok {
		return fmt.Errorf("invalid template data")
	}

	channels, err := json.Marshal(data.Channels)
	if err != nil {
		return err
	}

	template := &model.MessageTemplates{
		Code:            data.Code,
		Name:            data.Name,
		TitleTemplate:   data.TitleTemplate,
		ContentTemplate: data.ContentTemplate,
		Type:            int64(data.Type),
		Channels:        string(channels),
		Status:          1, // Default enabled
	}

	if data.Config != "" {
		template.Config = sql.NullString{
			String: data.Config,
			Valid:  true,
		}
	}

	_, err = h.templates.Insert(context.Background(), template)
	return err
}

func (h *TemplateHandler) handleTemplateUpdated(event *types.MessageEvent) error {
	data, ok := event.Data.(*types.TemplateData)
	if !ok {
		return fmt.Errorf("invalid template data")
	}

	template, err := h.templates.FindOneByCode(context.Background(), data.Code)
	if err != nil {
		return err
	}

	channels, err := json.Marshal(data.Channels)
	if err != nil {
		return err
	}

	template.Name = data.Name
	template.TitleTemplate = data.TitleTemplate
	template.ContentTemplate = data.ContentTemplate
	template.Channels = string(channels)

	if data.Config != "" {
		template.Config = sql.NullString{
			String: data.Config,
			Valid:  true,
		}
	}

	return h.templates.Update(context.Background(), template)
}