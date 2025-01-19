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
	userPointsFieldNames          = builder.RawFieldNames(&UserPoints{})
	userPointsRows                = strings.Join(userPointsFieldNames, ",")
	userPointsRowsExpectAutoSet   = strings.Join(stringx.Remove(userPointsFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	userPointsRowsWithPlaceHolder = strings.Join(stringx.Remove(userPointsFieldNames, "`user_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallMarketingUserPointsUserIdPrefix = "cache:mallMarketing:userPoints:userId:"
)

type (
	userPointsModel interface {
		Insert(ctx context.Context, data *UserPoints) (sql.Result, error)
		FindOne(ctx context.Context, userId uint64) (*UserPoints, error)
		Update(ctx context.Context, data *UserPoints) error
		Delete(ctx context.Context, userId uint64) error
	}

	defaultUserPointsModel struct {
		sqlc.CachedConn
		table string
	}

	UserPoints struct {
		UserId      uint64    `db:"user_id"`      // ç”¨æˆ·ID
		Points      int64     `db:"points"`       // ç§¯åˆ†ä½™é¢
		TotalPoints int64     `db:"total_points"` // ç´¯è®¡èŽ·å¾—ç§¯åˆ†
		UsedPoints  int64     `db:"used_points"`  // å·²ä½¿ç”¨ç§¯åˆ†
		UpdatedAt   time.Time `db:"updated_at"`   // æ›´æ–°æ—¶é—´
	}
)

func newUserPointsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultUserPointsModel {
	return &defaultUserPointsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`user_points`",
	}
}

func (m *defaultUserPointsModel) Delete(ctx context.Context, userId uint64) error {
	mallMarketingUserPointsUserIdKey := fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, userId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, userId)
	}, mallMarketingUserPointsUserIdKey)
	return err
}

func (m *defaultUserPointsModel) FindOne(ctx context.Context, userId uint64) (*UserPoints, error) {
	mallMarketingUserPointsUserIdKey := fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, userId)
	var resp UserPoints
	err := m.QueryRowCtx(ctx, &resp, mallMarketingUserPointsUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userPointsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, userId)
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

func (m *defaultUserPointsModel) Insert(ctx context.Context, data *UserPoints) (sql.Result, error) {
	mallMarketingUserPointsUserIdKey := fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, data.UserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, userPointsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.Points, data.TotalPoints, data.UsedPoints)
	}, mallMarketingUserPointsUserIdKey)
	return ret, err
}

func (m *defaultUserPointsModel) Update(ctx context.Context, data *UserPoints) error {
	mallMarketingUserPointsUserIdKey := fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, data.UserId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `user_id` = ?", m.table, userPointsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Points, data.TotalPoints, data.UsedPoints, data.UserId)
	}, mallMarketingUserPointsUserIdKey)
	return err
}

func (m *defaultUserPointsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, primary)
}

func (m *defaultUserPointsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userPointsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserPointsModel) tableName() string {
	return m.table
}
