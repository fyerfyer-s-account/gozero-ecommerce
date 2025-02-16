// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.5

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	couponsFieldNames          = builder.RawFieldNames(&Coupons{})
	couponsRows                = strings.Join(couponsFieldNames, ",")
	couponsRowsExpectAutoSet   = strings.Join(stringx.Remove(couponsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	couponsRowsWithPlaceHolder = strings.Join(stringx.Remove(couponsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallMarketingCouponsIdPrefix   = "cache:mallMarketing:coupons:id:"
	cacheMallMarketingCouponsCodePrefix = "cache:mallMarketing:coupons:code:"
)

type (
	couponsModel interface {
		Insert(ctx context.Context, data *Coupons) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Coupons, error)
		FindOneByCode(ctx context.Context, code string) (*Coupons, error)
		Update(ctx context.Context, data *Coupons) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultCouponsModel struct {
		sqlc.CachedConn
		table string
	}

	Coupons struct {
		Id        uint64       `db:"id"`         // ä¼˜æƒ åˆ¸ID
		Name      string       `db:"name"`       // ä¼˜æƒ åˆ¸åç§°
		Code      string       `db:"code"`       // ä¼˜æƒ åˆ¸ç 
		Type      int64        `db:"type"`       // ä¼˜æƒ åˆ¸ç±»åž‹ 1:æ»¡å‡ 2:æŠ˜æ‰£ 3:ç«‹å‡
		Value     float64      `db:"value"`      // ä¼˜æƒ é‡‘é¢æˆ–æŠ˜æ‰£çŽ‡
		MinAmount float64      `db:"min_amount"` // æœ€ä½Žä½¿ç”¨é‡‘é¢
		Status    int64        `db:"status"`     // çŠ¶æ€ 0:æœªå¼€å§‹ 1:è¿›è¡Œä¸­ 2:å·²ç»“æŸ 3:å·²å¤±æ•ˆ
		StartTime sql.NullTime `db:"start_time"` // å¼€å§‹æ—¶é—´
		EndTime   sql.NullTime `db:"end_time"`   // ç»“æŸæ—¶é—´
		Total     int64        `db:"total"`      // å‘è¡Œæ€»é‡
		Received  int64        `db:"received"`   // å·²é¢†å–æ•°é‡
		Used      int64        `db:"used"`       // å·²ä½¿ç”¨æ•°é‡
		PerLimit  int64        `db:"per_limit"`  // æ˜¯å¦é™åˆ¶æ¯äººé¢†å–æ•°é‡
		UserLimit int64        `db:"user_limit"` // æ¯äººé™é¢†æ•°é‡
		CreatedAt time.Time    `db:"created_at"` // åˆ›å»ºæ—¶é—´
		UpdatedAt time.Time    `db:"updated_at"` // æ›´æ–°æ—¶é—´
	}
)

func newCouponsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultCouponsModel {
	return &defaultCouponsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`coupons`",
	}
}

func (m *defaultCouponsModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallMarketingCouponsCodeKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsCodePrefix, data.Code)
	mallMarketingCouponsIdKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mallMarketingCouponsCodeKey, mallMarketingCouponsIdKey)
	return err
}

func (m *defaultCouponsModel) FindOne(ctx context.Context, id uint64) (*Coupons, error) {
	mallMarketingCouponsIdKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, id)
	var resp Coupons
	err := m.QueryRowCtx(ctx, &resp, mallMarketingCouponsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", couponsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCouponsModel) FindOneByCode(ctx context.Context, code string) (*Coupons, error) {
	mallMarketingCouponsCodeKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsCodePrefix, code)
	var resp Coupons
	err := m.QueryRowIndexCtx(ctx, &resp, mallMarketingCouponsCodeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `code` = ? limit 1", couponsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, code); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCouponsModel) Insert(ctx context.Context, data *Coupons) (sql.Result, error) {
	mallMarketingCouponsCodeKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsCodePrefix, data.Code)
	mallMarketingCouponsIdKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, couponsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Name, data.Code, data.Type, data.Value, data.MinAmount, data.Status, data.StartTime, data.EndTime, data.Total, data.Received, data.Used, data.PerLimit, data.UserLimit)
	}, mallMarketingCouponsCodeKey, mallMarketingCouponsIdKey)
	return ret, err
}

func (m *defaultCouponsModel) Update(ctx context.Context, newData *Coupons) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	mallMarketingCouponsCodeKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsCodePrefix, data.Code)
	mallMarketingCouponsIdKey := fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, couponsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.Name, newData.Code, newData.Type, newData.Value, newData.MinAmount, newData.Status, newData.StartTime, newData.EndTime, newData.Total, newData.Received, newData.Used, newData.PerLimit, newData.UserLimit, newData.Id)
	}, mallMarketingCouponsCodeKey, mallMarketingCouponsIdKey)
	return err
}

func (m *defaultCouponsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallMarketingCouponsIdPrefix, primary)
}

func (m *defaultCouponsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", couponsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultCouponsModel) tableName() string {
	return m.table
}
