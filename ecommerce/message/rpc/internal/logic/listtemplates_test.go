package logic

import (
	"context"
	"flag"
	"strings"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type ListTemplatesTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *ListTemplatesLogic
}

func TestListTemplatesSuite(t *testing.T) {
	suite.Run(t, new(ListTemplatesTestSuite))
}

func (s *ListTemplatesTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/message.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewListTemplatesLogic(context.Background(), s.ctx)
}

func (s *ListTemplatesTestSuite) SetupTest() {
	s.cleanData()
}

func (s *ListTemplatesTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *ListTemplatesTestSuite) cleanData() {
	// Clean test data using ListByTypeAndStatus
	templates, err := s.ctx.MessageTemplatesModel.ListByTypeAndStatus(context.Background(), 0, 0, 1, 100)
	if err != nil {
		return
	}
	
	for _, tpl := range templates {
		if strings.HasPrefix(tpl.Code, "TEST_TPL_") {
			_ = s.ctx.MessageTemplatesModel.Delete(context.Background(), tpl.Id)
		}
	}
}

func (s *ListTemplatesTestSuite) TestListTemplates() {
	tests := []struct {
		name    string
		setup   func() error
		req     *message.ListTemplatesRequest
		wantErr bool
		check   func(*message.ListTemplatesResponse)
	}{
		{
			name: "list all templates",
			setup: func() error {
				templates := []*model.MessageTemplates{
					{
						Code:            "TEST_TPL_1",
						Name:            "Test Template 1",
						TitleTemplate:   "Title 1",
						ContentTemplate: "Content 1",
						Type:            1,
						Channels:        "[1,2]",
						Status:          1,
					},
					{
						Code:            "TEST_TPL_2",
						Name:            "Test Template 2",
						TitleTemplate:   "Title 2",
						ContentTemplate: "Content 2",
						Type:            2,
						Channels:        "[2,3]",
						Status:          1,
					},
				}
				for _, tpl := range templates {
					_, err := s.ctx.MessageTemplatesModel.Insert(context.Background(), tpl)
					if err != nil {
						return err
					}
				}
				return nil
			},
			req: &message.ListTemplatesRequest{
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
			check: func(resp *message.ListTemplatesResponse) {
				assert.Equal(s.T(), int64(2), resp.Total)
				assert.Len(s.T(), resp.Templates, 2)
				if len(resp.Templates) >= 2 {
					assert.Equal(s.T(), "Test Template 2", resp.Templates[0].Name)
					assert.Equal(s.T(), "Test Template 1", resp.Templates[1].Name)
				}
			},
		},
		{
			name: "filter by type",
			setup: func() error {
				templates := []*model.MessageTemplates{
					{
						Code:            "TEST_TPL_1",
						Name:            "Test Template 1",
						TitleTemplate:   "Title 1",
						ContentTemplate: "Content 1",
						Type:            1,
						Channels:        "[1,2]",
						Status:          1,
					},
					{
						Code:            "TEST_TPL_2",
						Name:            "Test Template 2",
						TitleTemplate:   "Title 2",
						ContentTemplate: "Content 2",
						Type:            2,
						Channels:        "[2,3]",
						Status:          1,
					},
				}
				for _, tpl := range templates {
					_, err := s.ctx.MessageTemplatesModel.Insert(context.Background(), tpl)
					if err != nil {
						return err
					}
				}
				return nil
			},
			req: &message.ListTemplatesRequest{
				Type:     1,
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
			check: func(resp *message.ListTemplatesResponse) {
				assert.Equal(s.T(), int64(1), resp.Total)
				assert.Len(s.T(), resp.Templates, 1)
				assert.Equal(s.T(), "TEST_TPL_1", resp.Templates[0].Code)
			},
		},
		{
			name: "filter by status",
			setup: func() error {
				templates := []*model.MessageTemplates{
					{
						Code:            "TEST_TPL_1",
						Name:            "Test Template 1",
						TitleTemplate:   "Title 1",
						ContentTemplate: "Content 1",
						Type:            1,
						Channels:        "[1,2]",
						Status:          1,
					},
					{
						Code:            "TEST_TPL_2",
						Name:            "Test Template 2",
						TitleTemplate:   "Title 2",
						ContentTemplate: "Content 2",
						Type:            1,
						Channels:        "[2,3]",
						Status:          2,
					},
				}
				for _, tpl := range templates {
					_, err := s.ctx.MessageTemplatesModel.Insert(context.Background(), tpl)
					if err != nil {
						return err
					}
				}
				return nil
			},
			req: &message.ListTemplatesRequest{
				Status:   2,
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
			check: func(resp *message.ListTemplatesResponse) {
				assert.Equal(s.T(), int64(1), resp.Total)
				assert.Len(s.T(), resp.Templates, 1)
				assert.Equal(s.T(), "TEST_TPL_2", resp.Templates[0].Code)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.cleanData()
			if tt.setup != nil {
				err := tt.setup()
				assert.NoError(t, err)
			}

			resp, err := s.logic.ListTemplates(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if resp != nil && tt.check != nil {
				tt.check(resp)
			}
		})
	}
}
