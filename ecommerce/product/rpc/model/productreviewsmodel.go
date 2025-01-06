package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductReviewsModel = (*customProductReviewsModel)(nil)

type (
	// ProductReviewsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductReviewsModel.
	ProductReviewsModel interface {
		productReviewsModel
		FindManyByProductId(ctx context.Context, productId uint64, page, pageSize int) ([]*ProductReviews, error)
		FindOneByProductId(ctx context.Context, productId, userId uint64) (*ProductReviews, error)
		Count(ctx context.Context, productId uint64) (int64, error)
		BatchCreate(ctx context.Context, reviews []*ProductReviews) error
		UpdateReviews(ctx context.Context, id uint64, updates map[string]interface{}) error
	}

	customProductReviewsModel struct {
		*defaultProductReviewsModel
	}
)

// NewProductReviewsModel returns a model for the database table.
func NewProductReviewsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductReviewsModel {
	return &customProductReviewsModel{
		defaultProductReviewsModel: newProductReviewsModel(conn, c, opts...),
	}
}

func (m *customProductReviewsModel) FindManyByProductId(ctx context.Context, productId uint64, page, pageSize int) ([]*ProductReviews, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var reviews []*ProductReviews
	query := fmt.Sprintf("select %s from %s where `product_id` = ? limit ?, ?", productReviewsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &reviews, query, productId, offset, pageSize)

	return reviews, err
}

func (m *customProductReviewsModel) FindOneByProductId(ctx context.Context, productId, userId uint64) (*ProductReviews, error) {
	var review ProductReviews
	query := fmt.Sprintf("select %s from %s where `product_id` = ? and `user_id` = ?", productReviewsRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &review, query, productId, userId)

	return &review, err
}

func (m *customProductReviewsModel) Count(ctx context.Context, productId uint64) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where `product_id` = ?", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &count, query, productId)

	return count, err
}

func (m *customProductReviewsModel) UpdateStatus(ctx context.Context, id uint64, status int64) error {
	productReviewKey := fmt.Sprintf("%s%v", cacheMallProductProductReviewsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, productReviewKey)

	return err
}

func (m *customProductReviewsModel) UpdateContent(ctx context.Context, id uint64, rating int64, content, images sql.NullString) error {
	productReviewKey := fmt.Sprintf("%s%v", cacheMallProductProductReviewsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `rating` = ?, `content` = ?, `images` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, rating, content, images, id)
	}, productReviewKey)

	return err
}

func (m *customProductReviewsModel) BatchCreate(ctx context.Context, reviews []*ProductReviews) error {
	if len(reviews) == 0 {
		return nil
	}

	values := make([]string, 0, len(reviews))
	args := make([]interface{}, 0, len(reviews)*4)

	for _, review := range reviews {
		values = append(values, "(?, ?, ?, ?, ?)")
		args = append(args,
			review.ProductId,
			review.UserId,
			review.Rating,
			review.Content,
			review.Images)
	}

	query := fmt.Sprintf("insert into %s (product_id, user_id, rating, content, images) values %s",
		m.table, strings.Join(values, ","))

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	})

	return err
}

func (m *customProductReviewsModel) UpdateReviews(ctx context.Context, id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	var sets []string
	var args []interface{}
	for k, v := range updates {
		sets = append(sets, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	args = append(args, id)

	productReviewKey := fmt.Sprintf("%s%v", cacheMallProductProductReviewsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(sets, ", "))
		return conn.ExecCtx(ctx, query, args...)
	}, productReviewKey)

	return err
}