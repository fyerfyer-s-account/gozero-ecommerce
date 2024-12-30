package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoriesModel = (*customCategoriesModel)(nil)

var cacheMallProductCategoriesParentIdPrefix = "cache:mallProduct:categories:ParentId:"

type (
	// CategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoriesModel.
	CategoriesModel interface {
		categoriesModel
		FindByParentId(ctx context.Context, parentId uint64) ([]*Categories, error)
		FindOneByName(ctx context.Context, name string) (*Categories, error)
		UpdateSort(ctx context.Context, id uint64, sort int64) error
		GetLevel(ctx context.Context, parentId uint64) (int64, error)
		HasChildren(ctx context.Context, id uint64) (bool, error)
		HasProducts(ctx context.Context, id uint64) (bool, error)
	}

	customCategoriesModel struct {
		*defaultCategoriesModel
	}
)

// NewCategoriesModel returns a model for the database table.
func NewCategoriesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CategoriesModel {
	return &customCategoriesModel{
		defaultCategoriesModel: newCategoriesModel(conn, c, opts...),
	}
}

func (m *customCategoriesModel) FindByParentId(ctx context.Context, parentId uint64) ([]*Categories, error) {
	var resp []*Categories
	query := fmt.Sprintf("select %s from %s where `parent_id` = ? order by `sort`", categoriesRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, parentId)
	return resp, err
}

func (m *customCategoriesModel) FindOneByName(ctx context.Context, name string) (*Categories, error) {
	var resp Categories
	query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", categoriesRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, name)

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customCategoriesModel) UpdateSort(ctx context.Context, id uint64, sort int64) error {
	mallProductCategoriesIdKey := fmt.Sprintf("%s%v", cacheMallProductCategoriesIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `sort` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, sort, id)
	}, mallProductCategoriesIdKey)
	return err
}

func (m *customCategoriesModel) GetLevel(ctx context.Context, parentId uint64) (int64, error) {
	if parentId == 0 {
		return 1, nil
	}

	parent, err := m.FindOne(ctx, parentId)
	if err != nil {
		return 0, err
	}
	return parent.Level + 1, nil
}

func (m *customCategoriesModel) HasChildren(ctx context.Context, id uint64) (bool, error) {
	mallProductCategoriesParentIdKey := fmt.Sprintf("%s%v", cacheMallProductCategoriesParentIdPrefix, id)
	var count int64

	err := m.QueryRowCtx(ctx, &count, mallProductCategoriesParentIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select count(*) from %s where `parent_id` = ?", m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})

	return count > 0, err
}

func (m *customCategoriesModel) HasProducts(ctx context.Context, id uint64) (bool, error) {
	var count int64
	query := "select count(*) from products where `category_id` = ?"
	err := m.QueryRowNoCacheCtx(ctx, &count, query, id)
	return count > 0, err
}
