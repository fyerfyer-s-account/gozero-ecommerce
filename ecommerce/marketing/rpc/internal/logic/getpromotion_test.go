package logic

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/conf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetPromotionTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *GetPromotionLogic
}

func TestGetPromotionSuite(t *testing.T) {
    suite.Run(t, new(GetPromotionTestSuite))
}

func (s *GetPromotionTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
    s.logic = NewGetPromotionLogic(context.Background(), s.ctx)
}

func (s *GetPromotionTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *GetPromotionTestSuite) cleanData() {
    _ = s.ctx.PromotionsModel.Delete(context.Background(), 1001)
}

func (s *GetPromotionTestSuite) TestGetPromotion() {
    tests := []struct {
        name    string
        setup   func() error
        req     *marketing.GetPromotionRequest
        wantErr error
        check   func(*marketing.GetPromotionResponse)
    }{
        {
            name: "get existing promotion",
            setup: func() error {
                now := time.Now()
                promotion := &model.Promotions{
                    Name:      "Test Promotion",
                    Type:      1,
                    Rules:     "{\"minAmount\":100,\"discountAmount\":20}",
                    Status:    1,
                    StartTime: sql.NullTime{Time: now, Valid: true},
                    EndTime:   sql.NullTime{Time: now.Add(24 * time.Hour), Valid: true},
                }
                
                result, err := s.ctx.PromotionsModel.Insert(context.Background(), promotion)
                if err != nil {
                    return err
                }

                // Verify insertion
                affected, err := result.RowsAffected()
                if err != nil || affected == 0 {
                    return fmt.Errorf("failed to insert test promotion")
                }
                
                return nil
            },
            req: &marketing.GetPromotionRequest{
                Id: 1,
            },
            wantErr: nil,
            check: func(resp *marketing.GetPromotionResponse) {
                if resp == nil || resp.Promotion == nil {
                    s.T().Error("Response or promotion is nil")
                    return
                }
                assert.Equal(s.T(), "Test Promotion", resp.Promotion.Name)
                assert.Equal(s.T(), int32(1), resp.Promotion.Type)
                assert.Contains(s.T(), resp.Promotion.Rules, "minAmount")
                assert.Equal(s.T(), int32(1), resp.Promotion.Status)
            },
        },
        {
            name: "get non-existent promotion",
            req: &marketing.GetPromotionRequest{
                Id: 9999,
            },
            wantErr: zeroerr.ErrPromotionNotFound,
        },
        {
            name: "invalid promotion id",
            req: &marketing.GetPromotionRequest{
                Id: 0,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            // Clean before each test case
            s.cleanData()

            if tt.setup != nil {
                err := tt.setup()
                if !assert.NoError(t, err, "Setup failed") {
                    return
                }
            }

            resp, err := s.logic.GetPromotion(tt.req)
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