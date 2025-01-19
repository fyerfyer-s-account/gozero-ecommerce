package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserCouponsModel = (*customUserCouponsModel)(nil)

type (
	// UserCouponsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCouponsModel.
	UserCouponsModel interface {
		userCouponsModel
        Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
        FindByUserAndStatus(ctx context.Context, userId int64, status int32, page, pageSize int32) ([]*UserCoupons, error)
        CountByUser(ctx context.Context, userId int64, status int32) (int64, error)
        VerifyCoupon(ctx context.Context, userId, couponId int64) (*UserCoupons, error)
        UpdateStatus(ctx context.Context, id int64, status int32, orderNo string) error
        BatchInsert(ctx context.Context, coupons []*UserCoupons) error
        FindByOrderNo(ctx context.Context, orderNo string) (*UserCoupons, error)
        CountUserCoupon(ctx context.Context, userId, couponId int64) (int64, error)
	}

	customUserCouponsModel struct {
		*defaultUserCouponsModel
	}
)

// NewUserCouponsModel returns a model for the database table.
func NewUserCouponsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserCouponsModel {
	return &customUserCouponsModel{
		defaultUserCouponsModel: newUserCouponsModel(conn, c, opts...),
	}
}

func (m *customUserCouponsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
    return m.TransactCtx(ctx, fn)
}

func (m *customUserCouponsModel) FindByUserAndStatus(ctx context.Context, userId int64, status int32, page, pageSize int32) ([]*UserCoupons, error) {
    var coupons []*UserCoupons
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and `status` = ? order by id desc limit ?, ?", 
        userCouponsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &coupons, query, userId, status, (page-1)*pageSize, pageSize)
    return coupons, err
}

func (m *customUserCouponsModel) CountByUser(ctx context.Context, userId int64, status int32) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `status` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, userId, status)
    return count, err
}

func (m *customUserCouponsModel) VerifyCoupon(ctx context.Context, userId, couponId int64) (*UserCoupons, error) {
    var coupon UserCoupons
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and `coupon_id` = ? and `status` = 0 limit 1", 
        userCouponsRows, m.table)
    err := m.QueryRowNoCacheCtx(ctx, &coupon, query, userId, couponId)
    if err != nil {
        return nil, err
    }
    return &coupon, nil
}

func (m *customUserCouponsModel) UpdateStatus(ctx context.Context, id int64, status int32, orderNo string) error {
    key := fmt.Sprintf("%s%v", cacheMallMarketingUserCouponsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        query := fmt.Sprintf("update %s set `status` = ?, `used_time` = ?, `order_no` = ? where `id` = ?", m.table)
        return conn.ExecCtx(ctx, query, status, time.Now(), orderNo, id)
    }, key)
    return err
}

func (m *customUserCouponsModel) BatchInsert(ctx context.Context, coupons []*UserCoupons) error {
    if len(coupons) == 0 {
        return nil
    }
    
    values := make([]string, 0, len(coupons))
    args := make([]interface{}, 0, len(coupons)*5)
    for _, coupon := range coupons {
        values = append(values, "(?, ?, ?, ?, ?)")
        args = append(args, coupon.UserId, coupon.CouponId, coupon.Status, coupon.UsedTime, coupon.OrderNo)
    }
    
    query := fmt.Sprintf("insert into %s (%s) values %s", 
        m.table, userCouponsRowsExpectAutoSet, strings.Join(values, ","))
    _, err := m.ExecNoCacheCtx(ctx, query, args...)
    return err
}

func (m *customUserCouponsModel) FindByOrderNo(ctx context.Context, orderNo string) (*UserCoupons, error) {
    var coupon UserCoupons
    query := fmt.Sprintf("select %s from %s where `order_no` = ? limit 1", userCouponsRows, m.table)
    err := m.QueryRowNoCacheCtx(ctx, &coupon, query, orderNo)
    if err != nil {
        return nil, err
    }
    return &coupon, nil
}

func (m *customUserCouponsModel) CountUserCoupon(ctx context.Context, userId, couponId int64) (int64, error) {
    var count int64
    query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `coupon_id` = ?", m.table)
    err := m.QueryRowNoCacheCtx(ctx, &count, query, userId, couponId)
    return count, err
}