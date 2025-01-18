package logic

import (
    "context"
    "database/sql"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/message"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/zeromicro/go-zero/core/conf"
)

type GetTemplateTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *GetTemplateLogic
}

func TestGetTemplateSuite(t *testing.T) {
    suite.Run(t, new(GetTemplateTestSuite))
}

func (s *GetTemplateTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/message.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewGetTemplateLogic(context.Background(), s.ctx)
}

func (s *GetTemplateTestSuite) SetupTest() {
    s.cleanData()
}

func (s *GetTemplateTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *GetTemplateTestSuite) cleanData() {
    templates, _ := s.ctx.MessageTemplatesModel.FindByCode(context.Background(), "TEST_TPL")
    if templates != nil {
        _ = s.ctx.MessageTemplatesModel.Delete(context.Background(), templates.Id)
    }
}

func (s *GetTemplateTestSuite) TestGetTemplate() {
    tests := []struct {
        name    string
        setup   func() int64
        req     *message.GetTemplateRequest
        wantErr bool
        check   func(*message.GetTemplateResponse)
    }{
        {
            name: "normal case",
            setup: func() int64 {
                template := &model.MessageTemplates{
                    Code:            "TEST_TPL",
                    Name:            "Test Template",
                    TitleTemplate:   "Hello ${name}",
                    ContentTemplate: "Welcome to our platform, ${name}!",
                    Type:            1,
                    Channels:        "[1,2]",
                    Config: sql.NullString{
                        String: `{"key":"value"}`,
                        Valid:  true,
                    },
                    Status: 1,
                }
                result, _ := s.ctx.MessageTemplatesModel.Insert(context.Background(), template)
                id, _ := result.LastInsertId()
                return id
            },
            req: &message.GetTemplateRequest{
                Id: 1,
            },
            wantErr: false,
            check: func(resp *message.GetTemplateResponse) {
                assert.NotNil(s.T(), resp.Template)
                assert.Equal(s.T(), "TEST_TPL", resp.Template.Code)
                assert.Equal(s.T(), "Test Template", resp.Template.Name)
                assert.Equal(s.T(), "Hello ${name}", resp.Template.TitleTemplate)
                assert.Equal(s.T(), "Welcome to our platform, ${name}!", resp.Template.ContentTemplate)
                assert.Equal(s.T(), int32(1), resp.Template.Type)
                assert.ElementsMatch(s.T(), []int32{1, 2}, resp.Template.Channels)
                assert.JSONEq(s.T(), `{"key":"value"}`, resp.Template.Config)
                assert.Equal(s.T(), int32(1), resp.Template.Status)
                assert.Greater(s.T(), resp.Template.CreatedAt, int64(0))
                assert.Greater(s.T(), resp.Template.UpdatedAt, int64(0))
            },
        },
        {
            name:  "template not found",
            setup: func() int64 { return 0 },
            req: &message.GetTemplateRequest{
                Id: 999,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            msgId := tt.setup()
            if msgId > 0 {
                tt.req.Id = msgId
            }

            resp, err := s.logic.GetTemplate(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, resp)
            tt.check(resp)
        })
    }
}