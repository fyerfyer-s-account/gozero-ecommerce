package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserPointsModel = (*customUserPointsModel)(nil)

type (
	// UserPointsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserPointsModel.
	UserPointsModel interface {
		userPointsModel
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		IncrPoints(ctx context.Context, userId int64, points int64) error
		DecrPoints(ctx context.Context, userId int64, points int64) error
		Lock(ctx context.Context, session sqlx.Session, userId int64) error
		Unlock(ctx context.Context, session sqlx.Session, userId int64) error
		GetBalance(ctx context.Context, userId int64) (int64, error)
		InitUserPoints(ctx context.Context, userId int64) error
	}

	customUserPointsModel struct {
		*defaultUserPointsModel
	}
)

// NewUserPointsModel returns a model for the database table.
func NewUserPointsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserPointsModel {
	return &customUserPointsModel{
		defaultUserPointsModel: newUserPointsModel(conn, c, opts...),
	}
}

func (m *customUserPointsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.TransactCtx(ctx, fn)
}

func (m *customUserPointsModel) Lock(ctx context.Context, session sqlx.Session, userId int64) error {
	query := fmt.Sprintf("select user_id from %s where `user_id` = ? for update", m.table)
	var uid int64
	err := session.QueryRowCtx(ctx, &uid, query, userId)
	switch err {
	case nil:
		return nil
	case sqlx.ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}

// func (m *customUserPointsModel) InitUserPoints(ctx context.Context, userId int64) error {
//     // Use the provided session if within transaction
//     _, err := m.Insert(ctx, &UserPoints{
//         UserId:      uint64(userId),
//         Points:      0,
//         TotalPoints: 0,
//         UsedPoints:  0,
//         UpdatedAt:   time.Now(),
//     })
//     return err
// }

func (m *customUserPointsModel) Unlock(ctx context.Context, session sqlx.Session, userId int64) error {
	return nil
}

func (m *customUserPointsModel) IncrPoints(ctx context.Context, userId int64, points int64) error {
	key := fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, userId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set `points` = `points` + ?, `total_points` = `total_points` + ? where `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, points, points, userId)
	}, key)
	return err
}

func (m *customUserPointsModel) DecrPoints(ctx context.Context, userId int64, points int64) error {
	key := fmt.Sprintf("%s%v", cacheMallMarketingUserPointsUserIdPrefix, userId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set `points` = `points` - ?, `used_points` = `used_points` + ? where `user_id` = ? and `points` >= ?", m.table)
		return conn.ExecCtx(ctx, query, points, points, userId, points)
	}, key)
	return err
}

// func (m *customUserPointsModel) GetBalance(ctx context.Context, userId int64) (int64, error) {
//     userPoints, err := m.FindOne(ctx, uint64(userId))
//     if err != nil {
//         if err == ErrNotFound {
//             return 0, nil
//         }
//         return 0, err
//     }
//     return userPoints.Points, nil
// }

func (m *customUserPointsModel) GetBalance(ctx context.Context, userId int64) (int64, error) {
    points, err := m.FindOne(ctx, uint64(userId))
    if err == sqlx.ErrNotFound {
        return 0, zeroerr.ErrNotFound
    }
    if err != nil {
        return 0, err
    }
    return points.Points, nil
}

func (m *customUserPointsModel) InitUserPoints(ctx context.Context, userId int64) error {
    points := &UserPoints{
        UserId:      uint64(userId),
        Points:      0,
        TotalPoints: 0,
        UsedPoints:  0,
        UpdatedAt:   time.Now(),
    }
    _, err := m.Insert(ctx, points)
    return err
}
