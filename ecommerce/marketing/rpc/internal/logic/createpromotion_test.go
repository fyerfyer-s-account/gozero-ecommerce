package logic

import (
	"context"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/conf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CreatePromotionTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *CreatePromotionLogic
}

func TestCreatePromotionSuite(t *testing.T) {
    suite.Run(t, new(CreatePromotionTestSuite))
}

func (s *CreatePromotionTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
    s.logic = NewCreatePromotionLogic(context.Background(), s.ctx)
}

func (s *CreatePromotionTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *CreatePromotionTestSuite) cleanData() {
    _ = s.ctx.PromotionsModel.Delete(context.Background(), 1001)
    _ = s.ctx.PromotionsModel.Delete(context.Background(), 1002)
}

func (s *CreatePromotionTestSuite) TestCreatePromotion() {
    now := time.Now()
    tests := []struct {
        name    string
        req     *marketing.CreatePromotionRequest
        wantErr error
        check   func(*marketing.CreatePromotionResponse)
    }{
        {
            name: "valid fixed amount promotion",
            req: &marketing.CreatePromotionRequest{
                Name: "Test Promotion",
                Type: 1,
                Rules: mustMarshalRule(types.PromotionRule{
                    MinAmount:      100,
                    DiscountAmount: 20,
                }),
                StartTime: now.Unix(),
                EndTime:   now.Add(24 * time.Hour).Unix(),
            },
            wantErr: nil,
            check: func(resp *marketing.CreatePromotionResponse) {
                assert.NotZero(s.T(), resp.Id)

                // Verify database
                promo, err := s.ctx.PromotionsModel.FindOne(context.Background(), uint64(resp.Id))
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), "Test Promotion", promo.Name)
                assert.Equal(s.T(), int64(1), promo.Type)
            },
        },
        {
            name: "invalid promotion type",
            req: &marketing.CreatePromotionRequest{
                Name: "Invalid Type",
                Type: 4,
                Rules: mustMarshalRule(types.PromotionRule{
                    MinAmount:      100,
                    DiscountAmount: 20,
                }),
                StartTime: now.Unix(),
                EndTime:   now.Add(24 * time.Hour).Unix(),
            },
            wantErr: zeroerr.ErrInvalidPromotionType,
        },
        {
            name: "expired end time",
            req: &marketing.CreatePromotionRequest{
                Name: "Expired Promotion",
                Type: 1,
                Rules: mustMarshalRule(types.PromotionRule{
                    MinAmount:      100,
                    DiscountAmount: 20,
                }),
                StartTime: now.Add(-48 * time.Hour).Unix(),
                EndTime:   now.Add(-24 * time.Hour).Unix(),
            },
            wantErr: zeroerr.ErrMarketingExpired,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()

            resp, err := s.logic.CreatePromotion(tt.req)
            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                return
            }

            assert.NoError(t, err)
            if tt.check != nil {
                tt.check(resp)
            }
        })
    }
}

func mustMarshalRule(rule types.PromotionRule) string {
    data, err := json.Marshal(rule)
    if err != nil {
        panic(err)
    }
    return string(data)
}