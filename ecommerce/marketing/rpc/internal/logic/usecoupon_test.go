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
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UseCouponTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *UseCouponLogic
}

func TestUseCouponSuite(t *testing.T) {
    suite.Run(t, new(UseCouponTestSuite))
}

func (s *UseCouponTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/marketing.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewUseCouponLogic(context.Background(), s.ctx)
}

func (s *UseCouponTestSuite) TearDownTest() {
    if s.ctx != nil && s.ctx.UserCouponsModel != nil {
        ctx := context.Background()
        _ = s.ctx.UserCouponsModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
            // Clean up test data
            _, _ = session.ExecCtx(ctx, "DELETE FROM user_coupons WHERE user_id = ?", 1001)
            _, _ = session.ExecCtx(ctx, "DELETE FROM coupons WHERE id IN (SELECT coupon_id FROM user_coupons WHERE user_id = ?)", 1001)
            return nil
        })
    }
}

func (s *UseCouponTestSuite) TestUseCoupon() {
    tests := []struct {
        name    string
        setup   func() (int64, error)
        req     *marketing.UseCouponRequest
        wantErr error
        check   func(int64, *marketing.UseCouponResponse)
    }{
        {
            name: "use valid coupon",
            setup: func() (int64, error) {
                ctx := context.Background()
                
                // Create test coupon
                coupon := &model.Coupons{
                    Name:      "Test Coupon",
                    Code:      "TEST" + time.Now().Format("150405"),
                    Type:      1,
                    Value:     20,
                    MinAmount: 100,
                    Status:    1,
                    StartTime: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
                    EndTime:   sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
                    Total:     100,
                    Received:  1,
                    Used:      0,
                }
                
                result, err := s.ctx.CouponsModel.Insert(ctx, coupon)
                if err != nil {
                    return 0, err
                }
                
                couponId, err := result.LastInsertId()
                if err != nil {
                    return 0, err
                }

                // Create user coupon
                userCoupon := &model.UserCoupons{
                    UserId:   1001,
                    CouponId: uint64(couponId),
                    Status:   0,
                }
                _, err = s.ctx.UserCouponsModel.Insert(ctx, userCoupon)
                return couponId, err
            },
            req: &marketing.UseCouponRequest{
                UserId:   1001,
                OrderNo:  "TEST123",
            },
            check: func(couponId int64, resp *marketing.UseCouponResponse) {
                assert.True(s.T(), resp.Success)
                
                // Verify coupon used count
                coupon, err := s.ctx.CouponsModel.FindOne(context.Background(), uint64(couponId))
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), int64(1), coupon.Used)

                // Verify user coupon status
                userCoupons, err := s.ctx.UserCouponsModel.FindByUserAndStatus(context.Background(), 1001, 1, 1, 10)
                assert.NoError(s.T(), err)
                assert.Len(s.T(), userCoupons, 1)
                assert.Equal(s.T(), "TEST123", userCoupons[0].OrderNo.String)
            },
        },
        {
            name: "invalid parameters",
            req: &marketing.UseCouponRequest{
                UserId:   0,
                CouponId: 1,
                OrderNo:  "",
            },
            wantErr: zeroerr.ErrInvalidMarketingParam,
        },
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            // Clean up before each test
            s.TearDownTest()
            
            var couponId int64
            if tt.setup != nil {
                var err error
                couponId, err = tt.setup()
                s.NoError(err)
                if tt.req.CouponId == 0 {
                    tt.req.CouponId = couponId
                }
            }

            resp, err := s.logic.UseCoupon(tt.req)
            if tt.wantErr != nil {
                s.Equal(tt.wantErr, err)
                return
            }

            s.NoError(err)
            if tt.check != nil {
                tt.check(couponId, resp)
            }
        })
    }
}