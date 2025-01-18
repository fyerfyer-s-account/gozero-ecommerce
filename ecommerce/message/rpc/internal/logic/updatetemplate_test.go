package logic

import (
	"context"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type UpdateTemplateTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *UpdateTemplateLogic
}

func TestUpdateTemplateSuite(t *testing.T) {
	suite.Run(t, new(UpdateTemplateTestSuite))
}

func (s *UpdateTemplateTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/message.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewUpdateTemplateLogic(context.Background(), s.ctx)
}

func (s *UpdateTemplateTestSuite) SetupTest() {
	s.cleanData()
}

func (s *UpdateTemplateTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *UpdateTemplateTestSuite) cleanData() {
	templates, _ := s.ctx.MessageTemplatesModel.FindByCode(context.Background(), "TEST_TPL")
	if templates != nil {
		_ = s.ctx.MessageTemplatesModel.Delete(context.Background(), templates.Id)
	}
}

func (s *UpdateTemplateTestSuite) TestUpdateTemplate() {
	tests := []struct {
		name    string
		setup   func() uint64
		req     *message.UpdateTemplateRequest
		wantErr error
	}{
		{
			name: "normal update",
			setup: func() uint64 {
				template := &model.MessageTemplates{
					Code:            "TEST_TPL",
					Name:            "Test Template",
					TitleTemplate:   "Hello ${name}",
					ContentTemplate: "Welcome ${name}",
					Type:            1,
					Channels:        "[1,2]",
					Status:          1,
				}
				result, _ := s.ctx.MessageTemplatesModel.Insert(context.Background(), template)
				id, _ := result.LastInsertId()
				return uint64(id)
			},
			req: &message.UpdateTemplateRequest{
				Name:            "Updated Template",
				TitleTemplate:   "Hi ${name}",
				ContentTemplate: "Welcome back ${name}",
				Channels:        []int32{1, 2, 3},
				Config:          `{"key":"value"}`,
				Status:          2,
			},
			wantErr: nil,
		},
		{
			name:  "template not found",
			setup: func() uint64 { return 0 },
			req: &message.UpdateTemplateRequest{
				Id:              999,
				Name:            "Test",
				TitleTemplate:   "Test",
				ContentTemplate: "Test",
			},
			wantErr: zeroerr.ErrTemplateNotFound,
		},
		{
			name: "invalid status",
			setup: func() uint64 {
				template := &model.MessageTemplates{
					Code:            "TEST_TPL",
					Name:            "Test Template",
					TitleTemplate:   "Hello ${name}",
					ContentTemplate: "Welcome ${name}",
					Type:            1,
					Channels:        "[1,2]",
					Status:          1,
				}
				result, _ := s.ctx.MessageTemplatesModel.Insert(context.Background(), template)
				id, _ := result.LastInsertId()
				return uint64(id)
			},
			req: &message.UpdateTemplateRequest{
				Name:            "Test",
				TitleTemplate:   "Test",
				ContentTemplate: "Test",
				Status:          3, // Invalid status
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.cleanData()
			id := tt.setup()
			if id > 0 {
				tt.req.Id = int64(id)
			}

			resp, err := s.logic.UpdateTemplate(tt.req)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, resp.Success)

			// Verify updated template
			template, err := s.ctx.MessageTemplatesModel.FindOne(context.Background(), id)
			assert.NoError(t, err)
			assert.Equal(t, tt.req.Name, template.Name)
			assert.Equal(t, tt.req.TitleTemplate, template.TitleTemplate)
			assert.Equal(t, tt.req.ContentTemplate, template.ContentTemplate)
			assert.Equal(t, int64(tt.req.Status), template.Status)
			assert.JSONEq(t, tt.req.Config, template.Config.String)
		})
	}
}
