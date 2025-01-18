package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type CreateTemplateTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *CreateTemplateLogic
}

func TestCreateTemplateSuite(t *testing.T) {
    suite.Run(t, new(CreateTemplateTestSuite))
}

func (s *CreateTemplateTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/message.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewCreateTemplateLogic(context.Background(), s.ctx)
}

func (s *CreateTemplateTestSuite) SetupTest() {
    s.cleanData()
}

func (s *CreateTemplateTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *CreateTemplateTestSuite) cleanData() {
    // Clean test data
    templates, _ := s.ctx.MessageTemplatesModel.FindByCode(context.Background(), "TEST_TPL")
    if templates != nil {
        _ = s.ctx.MessageTemplatesModel.Delete(context.Background(), templates.Id)
    }
}

func (s *CreateTemplateTestSuite) TestCreateTemplate() {
    tests := []struct {
        name    string
        req     *message.CreateTemplateRequest
        wantErr error
        check   func(*message.CreateTemplateResponse)
    }{
        {
            name: "normal case",
            req: &message.CreateTemplateRequest{
                Code:            "TEST_TPL",
                Name:            "Test Template",
                TitleTemplate:   "Hello ${name}",
                ContentTemplate: "Welcome to our platform, ${name}!",
                Type:           1,
                Channels:       []int32{1, 2},
                Config:         `{"key":"value"}`,
            },
            wantErr: nil,
            check: func(resp *message.CreateTemplateResponse) {
                assert.Greater(s.T(), resp.TemplateId, int64(0))

                // Verify in database
                template, err := s.ctx.MessageTemplatesModel.FindOne(context.Background(), uint64(resp.TemplateId))
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), "TEST_TPL", template.Code)
                assert.Equal(s.T(), "Test Template", template.Name)
                assert.Equal(s.T(), "[1,2]", template.Channels)
                assert.Equal(s.T(), sql.NullString{String: `{"key":"value"}`, Valid: true}, template.Config)
            },
        },
        {
            name: "missing required fields",
            req: &message.CreateTemplateRequest{
                Code: "TEST_TPL",
            },
            wantErr: zeroerr.ErrInvalidTemplate,
        },
        {
            name: "duplicate code",
            req: &message.CreateTemplateRequest{
                Code:            "TEST_TPL",
                Name:            "Test Template",
                TitleTemplate:   "Hello ${name}",
                ContentTemplate: "Welcome to our platform, ${name}!",
                Type:           1,
                Channels:       []int32{1, 2},
            },
            wantErr: zeroerr.ErrDuplicateTemplate,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.CreateTemplate(tt.req)
            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                return
            }
            assert.NoError(t, err)
            tt.check(resp)
        })
    }
}
