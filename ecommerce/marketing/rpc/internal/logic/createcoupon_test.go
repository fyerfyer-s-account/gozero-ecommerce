package logic

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/conf"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CreateCouponTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *CreateCouponLogic
}

func TestCreateCouponSuite(t *testing.T) {
    suite.Run(t, new(CreateCouponTestSuite))
}

func (s *CreateCouponTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
    s.logic = NewCreateCouponLogic(context.Background(), s.ctx)
}

func (s *CreateCouponTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *CreateCouponTestSuite) cleanData() {
    _ = s.ctx.CouponsModel.Delete(context.Background(), 1001)
    _ = s.ctx.CouponsModel.Delete(context.Background(), 1002)
}

func (s *CreateCouponTestSuite) TestCreateCoupon() {
    now := time.Now()
    tests := []struct {
        name    string
        req     *marketing.CreateCouponRequest
        wantErr error
        check   func(*marketing.CreateCouponResponse)
    }{
        {
            name: "valid fixed amount coupon",
            req: &marketing.CreateCouponRequest{
                Name:      "Test Coupon",
                Type:      1,
                Value:     20,
                MinAmount: 100,
                StartTime: now.Unix(),
                EndTime:   now.Add(24 * time.Hour).Unix(),
                Total:     1000,
                PerLimit:  1,
                UserLimit: 1,
            },
            wantErr: nil,
            check: func(resp *marketing.CreateCouponResponse) {
                assert.NotZero(s.T(), resp.Id)
                assert.NotEmpty(s.T(), resp.Code)

                // Verify database
                coupon, err := s.ctx.CouponsModel.FindOne(context.Background(), uint64(resp.Id))
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), "Test Coupon", coupon.Name)
                assert.Equal(s.T(), int64(1), coupon.Type)
                assert.Equal(s.T(), 20.0, coupon.Value)
                assert.Equal(s.T(), 100.0, coupon.MinAmount)
            },
        },
        {
            name: "invalid coupon type",
            req: &marketing.CreateCouponRequest{
                Name:      "Invalid Type",
                Type:      4,
                Value:     20,
                MinAmount: 100,
                StartTime: now.Unix(),
                EndTime:   now.Add(24 * time.Hour).Unix(),
                Total:     1000,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
        {
            name: "expired end time",
            req: &marketing.CreateCouponRequest{
                Name:      "Expired Coupon",
                Type:      1,
                Value:     20,
                MinAmount: 100,
                StartTime: now.Add(-48 * time.Hour).Unix(),
                EndTime:   now.Add(-24 * time.Hour).Unix(),
                Total:     1000,
            },
            wantErr: zeroerr.ErrMarketingExpired,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()

            resp, err := s.logic.CreateCoupon(tt.req)
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