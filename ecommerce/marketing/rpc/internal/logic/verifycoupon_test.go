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

type VerifyCouponTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *VerifyCouponLogic
}

func TestVerifyCouponSuite(t *testing.T) {
    suite.Run(t, new(VerifyCouponTestSuite))
}

func (s *VerifyCouponTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewVerifyCouponLogic(context.Background(), s.ctx)
}

func (s *VerifyCouponTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *VerifyCouponTestSuite) cleanData() {
    _ = s.ctx.CouponsModel.Delete(context.Background(), 1001)
    _ = s.ctx.UserCouponsModel.Delete(context.Background(), 1001)
}

func (s *VerifyCouponTestSuite) TestVerifyCoupon() {
    tests := []struct {
        name    string
        setup   func() (int64, error)
        req     *marketing.VerifyCouponRequest
        wantErr error
        check   func(*marketing.VerifyCouponResponse)
    }{
        {
            name: "verify valid fixed amount coupon",
            setup: func() (int64, error) {
                now := time.Now()
                // Create unique code for each test
                code := fmt.Sprintf("TEST%s%d", time.Now().Format("150405"), time.Now().UnixNano())
                result, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                    Name:      "Test Coupon",
                    Code:      code,
                    Type:      1, // Fixed amount
                    Value:     20,
                    MinAmount: 100,
                    Status:    1,
                    StartTime: sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
                    EndTime:   sql.NullTime{Time: now.Add(time.Hour), Valid: true},
                })
                if err != nil {
                    return 0, err
                }

                couponId, err := result.LastInsertId()
                if err != nil {
                    return 0, err
                }

                // Create user coupon
                _, err = s.ctx.UserCouponsModel.Insert(context.Background(), &model.UserCoupons{
                    UserId:   1001,
                    CouponId: uint64(couponId),
                    Status:   0,
                })
                return couponId, err
            },
            req: &marketing.VerifyCouponRequest{
                UserId:   1001,
                Amount:   200,
            },
            check: func(resp *marketing.VerifyCouponResponse) {
                if !assert.NotNil(s.T(), resp) {
                    return
                }
                assert.True(s.T(), resp.Valid)
                assert.Equal(s.T(), float64(20), resp.DiscountAmount)
            },
        },
        {
            name: "insufficient order amount",
            setup: func() (int64, error) {
                now := time.Now()
                code := fmt.Sprintf("TEST%s%d", time.Now().Format("150405"), time.Now().UnixNano())
                result, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                    Name:      "Test Coupon",
                    Code:      code,
                    Type:      1,
                    Value:     20,
                    MinAmount: 100,
                    Status:    1,
                    StartTime: sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
                    EndTime:   sql.NullTime{Time: now.Add(time.Hour), Valid: true},
                })
                if err != nil {
                    return 0, err
                }

                couponId, err := result.LastInsertId()
                if err != nil {
                    return 0, err
                }

                _, err = s.ctx.UserCouponsModel.Insert(context.Background(), &model.UserCoupons{
                    UserId:   1001,
                    CouponId: uint64(couponId),
                    Status:   0,
                })
                return couponId, err
            },
            req: &marketing.VerifyCouponRequest{
                UserId:   1001,
                Amount:   50,
            },
            check: func(resp *marketing.VerifyCouponResponse) {
                if !assert.NotNil(s.T(), resp) {
                    return
                }
                assert.False(s.T(), resp.Valid)
                assert.Contains(s.T(), resp.Message, "must be greater than")
            },
        },
        {
            name: "invalid parameters",
            req: &marketing.VerifyCouponRequest{
                UserId:   0,
                CouponId: 0,
                Amount:   0,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            var couponId int64
            if tt.setup != nil {
                var err error
                couponId, err = tt.setup()
                if !assert.NoError(t, err) {
                    return
                }
                if tt.req.CouponId == 0 {
                    tt.req.CouponId = couponId
                }
            }

            resp, err := s.logic.VerifyCoupon(tt.req)
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