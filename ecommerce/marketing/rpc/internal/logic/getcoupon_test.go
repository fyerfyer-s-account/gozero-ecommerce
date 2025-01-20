package logic

import (
	"context"
	"database/sql"
	"flag"
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

type GetCouponTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *GetCouponLogic
}

func TestGetCouponSuite(t *testing.T) {
    suite.Run(t, new(GetCouponTestSuite))
}

func (s *GetCouponTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
    s.logic = NewGetCouponLogic(context.Background(), s.ctx)
}

func (s *GetCouponTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *GetCouponTestSuite) cleanData() {
    rows, err := s.ctx.CouponsModel.FindManyByIds(context.Background(), []uint64{1001, 1002})
    if err == nil {
        for _, row := range rows {
            _ = s.ctx.CouponsModel.Delete(context.Background(), row.Id)
        }
    }
}

func (s *GetCouponTestSuite) TestGetCoupon() {
    tests := []struct {
        name    string
        setup   func() (uint64, error) // Changed to return ID
        req     *marketing.GetCouponRequest
        wantErr error
        check   func(*marketing.GetCouponResponse)
    }{
        {
            name: "get existing coupon",
            setup: func() (uint64, error) {
                now := time.Now()
                result, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                    Name:      "Test Coupon",
                    Code:      "TEST123",
                    Type:      1,
                    Value:     20,
                    MinAmount: 100,
                    Status:    1,
                    StartTime: sql.NullTime{Time: now, Valid: true},
                    EndTime:   sql.NullTime{Time: now.Add(24 * time.Hour), Valid: true},
                    Total:     1000,
                    Received:  10,
                    Used:      5,
                    PerLimit:  1,
                    UserLimit: 1,
                })
                if err != nil {
                    return 0, err
                }
                id, err := result.LastInsertId()
                return uint64(id), err
            },
            req:     &marketing.GetCouponRequest{},
            wantErr: nil,
            check: func(resp *marketing.GetCouponResponse) {
                assert.NotNil(s.T(), resp.Coupon)
                assert.Equal(s.T(), "Test Coupon", resp.Coupon.Name)
                // Remove code check since it's auto-generated
                assert.Equal(s.T(), int32(1), resp.Coupon.Type)
                assert.Equal(s.T(), float64(20), resp.Coupon.Value)
                assert.Equal(s.T(), float64(100), resp.Coupon.MinAmount)
                assert.Equal(s.T(), int32(1), resp.Coupon.Status)
                // Add length check for code
                assert.Len(s.T(), resp.Coupon.Code, 12) // Assuming code length is 12
            },
        },
        {
            name: "get non-existent coupon",
            req: &marketing.GetCouponRequest{
                Id: 9999,
            },
            wantErr: zeroerr.ErrCouponNotFound,
        },
        {
            name: "invalid coupon id",
            req: &marketing.GetCouponRequest{
                Id: 0,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            var id uint64
            if tt.setup != nil {
                var err error
                id, err = tt.setup()
                assert.NoError(t, err)
                if tt.req.Id == 0 { 
                    tt.req.Id = int64(id)
                }
            }

            resp, err := s.logic.GetCoupon(tt.req)
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