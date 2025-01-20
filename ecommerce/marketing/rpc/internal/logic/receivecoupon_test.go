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

type ReceiveCouponTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *ReceiveCouponLogic
}

func TestReceiveCouponSuite(t *testing.T) {
    suite.Run(t, new(ReceiveCouponTestSuite))
}

func (s *ReceiveCouponTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewReceiveCouponLogic(context.Background(), s.ctx)
}

func (s *ReceiveCouponTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *ReceiveCouponTestSuite) cleanData() {
    // Clean user coupons by test user ID
    userCoupons, err := s.ctx.UserCouponsModel.FindByUserAndStatus(context.Background(), 1001, 0, 1, 100)
    if err == nil {
        for _, uc := range userCoupons {
            _ = s.ctx.UserCouponsModel.Delete(context.Background(), uc.Id)
        }
    }
}

func (s *ReceiveCouponTestSuite) TestReceiveCoupon() {
    tests := []struct {
        name    string
        setup   func() (int64, error)
        req     *marketing.ReceiveCouponRequest
        wantErr error
        check   func(int64, *marketing.ReceiveCouponResponse)
    }{
        {
            name: "receive available coupon",
            setup: func() (int64, error) {
                now := time.Now()
                // Generate unique code for test
                code := fmt.Sprintf("TEST%d", time.Now().UnixNano())
                result, err := s.ctx.CouponsModel.Insert(context.Background(), &model.Coupons{
                    Name:      "Test Coupon",
                    Code:      code,
                    Type:      1,
                    Value:     20,
                    MinAmount: 100,
                    Status:    1,
                    StartTime: sql.NullTime{Time: now.Add(-time.Hour), Valid: true},
                    EndTime:   sql.NullTime{Time: now.Add(time.Hour), Valid: true},
                    Total:     100,
                    Received:  0,
                    Used:      0,
                    PerLimit:  1,
                    UserLimit: 1,
                })
                if err != nil {
                    return 0, err
                }
                return result.LastInsertId()
            },
            req: &marketing.ReceiveCouponRequest{
                UserId: 1001,
            },
            check: func(couponId int64, resp *marketing.ReceiveCouponResponse) {
                assert.True(s.T(), resp.Success)

                // Verify coupon received count
                coupon, err := s.ctx.CouponsModel.FindOne(context.Background(), uint64(couponId))
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), int64(1), coupon.Received)

                // Verify user coupon exists
                userCoupons, err := s.ctx.UserCouponsModel.FindByUserAndStatus(context.Background(), 1001, 0, 1, 10)
                assert.NoError(s.T(), err)
                assert.NotEmpty(s.T(), userCoupons)
                assert.Equal(s.T(), uint64(couponId), userCoupons[0].CouponId)
            },
        },
        {
            name: "invalid user id",
            req: &marketing.ReceiveCouponRequest{
                UserId:   0,
                CouponId: 1,
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
        {
            name: "non-existent coupon",
            req: &marketing.ReceiveCouponRequest{
                UserId:   1001,
                CouponId: 9999,
            },
            wantErr: zeroerr.ErrCouponNotFound,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            s.cleanData()
            var couponId int64
            if tt.setup != nil {
                var err error
                couponId, err = tt.setup()
                assert.NoError(t, err)
                if tt.req.CouponId == 0 {
                    tt.req.CouponId = couponId
                }
            }

            resp, err := s.logic.ReceiveCoupon(tt.req)
            if tt.wantErr != nil {
                assert.Equal(t, tt.wantErr, err)
                return
            }

            assert.NoError(t, err)
            if tt.check != nil {
                tt.check(couponId, resp)
            }
        })
    }
}