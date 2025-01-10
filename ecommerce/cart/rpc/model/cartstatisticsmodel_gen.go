// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.4

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
	cartStatisticsFieldNames          = builder.RawFieldNames(&CartStatistics{})
	cartStatisticsRows                = strings.Join(cartStatisticsFieldNames, ",")
	cartStatisticsRowsExpectAutoSet   = strings.Join(stringx.Remove(cartStatisticsFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	cartStatisticsRowsWithPlaceHolder = strings.Join(stringx.Remove(cartStatisticsFieldNames, "`user_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallCartCartStatisticsUserIdPrefix = "cache:mallCart:cartStatistics:userId:"
)

type (
	cartStatisticsModel interface {
		Insert(ctx context.Context, data *CartStatistics) (sql.Result, error)
		FindOne(ctx context.Context, userId uint64) (*CartStatistics, error)
		Update(ctx context.Context, data *CartStatistics) error
		Delete(ctx context.Context, userId uint64) error
	}

	defaultCartStatisticsModel struct {
		sqlc.CachedConn
		table string
	}

	CartStatistics struct {
		UserId           uint64    `db:"user_id"`           // ç”¨æˆ·ID
		TotalQuantity    int64     `db:"total_quantity"`    // å•†å“æ€»æ•°é‡
		SelectedQuantity int64     `db:"selected_quantity"` // å·²é€‰å•†å“æ•°é‡
		TotalAmount      float64   `db:"total_amount"`      // å•†å“æ€»é‡‘é¢
		SelectedAmount   float64   `db:"selected_amount"`   // å·²é€‰å•†å“é‡‘é¢
		UpdatedAt        time.Time `db:"updated_at"`        // æ›´æ–°æ—¶é—´
	}
)

func newCartStatisticsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultCartStatisticsModel {
	return &defaultCartStatisticsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`cart_statistics`",
	}
}

func (m *defaultCartStatisticsModel) Delete(ctx context.Context, userId uint64) error {
	mallCartCartStatisticsUserIdKey := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, userId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, userId)
	}, mallCartCartStatisticsUserIdKey)
	return err
}

func (m *defaultCartStatisticsModel) FindOne(ctx context.Context, userId uint64) (*CartStatistics, error) {
	mallCartCartStatisticsUserIdKey := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, userId)
	var resp CartStatistics
	err := m.QueryRowCtx(ctx, &resp, mallCartCartStatisticsUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", cartStatisticsRows, m.table)
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

func (m *defaultCartStatisticsModel) Insert(ctx context.Context, data *CartStatistics) (sql.Result, error) {
	mallCartCartStatisticsUserIdKey := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, data.UserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, cartStatisticsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.TotalQuantity, data.SelectedQuantity, data.TotalAmount, data.SelectedAmount)
	}, mallCartCartStatisticsUserIdKey)
	return ret, err
}

func (m *defaultCartStatisticsModel) Update(ctx context.Context, data *CartStatistics) error {
	mallCartCartStatisticsUserIdKey := fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, data.UserId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `user_id` = ?", m.table, cartStatisticsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.TotalQuantity, data.SelectedQuantity, data.TotalAmount, data.SelectedAmount, data.UserId)
	}, mallCartCartStatisticsUserIdKey)
	return err
}

func (m *defaultCartStatisticsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallCartCartStatisticsUserIdPrefix, primary)
}

func (m *defaultCartStatisticsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", cartStatisticsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultCartStatisticsModel) tableName() string {
	return m.table
}
