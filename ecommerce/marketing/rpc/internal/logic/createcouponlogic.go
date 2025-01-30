package logic

import (
    "context"
    "database/sql"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/marketing"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/util"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateCouponLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewCreateCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCouponLogic {
    return &CreateCouponLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// 优惠券管理
func (l *CreateCouponLogic) CreateCoupon(in *marketing.CreateCouponRequest) (*marketing.CreateCouponResponse, error) {
    // Validate input
    if err := l.validateInput(in); err != nil {
        return nil, err
    }
    
    var code string 
    var couponId uint64
    err := l.svcCtx.CouponsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        // Generate unique code using util package
        code = util.GenerateRandomString(12)

        // Create coupon
        result, err := l.svcCtx.CouponsModel.Insert(ctx, &model.Coupons{
            Name:      in.Name,
            Code:      code,
            Type:      int64(in.Type),
            Value:     in.Value,
            MinAmount: in.MinAmount,
            Status:    1, // Active
            StartTime: sql.NullTime{Time: time.Unix(in.StartTime, 0), Valid: true},
            EndTime:   sql.NullTime{Time: time.Unix(in.EndTime, 0), Valid: true},
            Total:     int64(in.Total),
            Received:  0,
            Used:      0,
            PerLimit:  in.PerLimit,
            UserLimit: int64(in.UserLimit),
        })
        if err != nil {
            return err
        }

        id, err := result.LastInsertId()
        if err != nil {
            return err
        }
        couponId = uint64(id)

        return nil
    })

    if err != nil {
        return nil, err
    }

    return &marketing.CreateCouponResponse{
        Id:   int64(couponId),
        Code: code,
    }, nil
}

func (l *CreateCouponLogic) validateInput(in *marketing.CreateCouponRequest) error {
    if in.Name == "" {
        return zeroerr.ErrInvalidMarketingParam
    }
    if in.Type < 1 || in.Type > 3 {
        return zeroerr.ErrInvalidMarketingParam
    }
    if in.Value <= 0 {
        return zeroerr.ErrInvalidMarketingParam
    }
    if in.Total <= 0 {
        return zeroerr.ErrInvalidMarketingParam
    }
    if in.StartTime >= in.EndTime {
        return zeroerr.ErrInvalidMarketingParam
    }
    if time.Unix(in.EndTime, 0).Before(time.Now()) {
        return zeroerr.ErrMarketingExpired
    }
    return nil
}