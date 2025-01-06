package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ ProductsModel = (*customProductsModel)(nil)

type (
	ProductsModel interface {
		productsModel
		FindManyByCategoryId(ctx context.Context, categoryId uint64, page, pageSize int) ([]*Products, error)
		FindOneByName(ctx context.Context, name string) (*Products, error)
		SearchByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*Products, error)
		GeneralSearch(ctx context.Context, page, pageSize int) ([]*Products, error)
		Count(ctx context.Context, categoryId uint64, keyword string) (int64, error)
		UpdatePartial(ctx context.Context, id uint64, updates map[string]interface{}) error
	}

	customProductsModel struct {
		*defaultProductsModel
	}
)

// NewProductsModel returns a model for the database table.
func NewProductsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductsModel {
	return &customProductsModel{
		defaultProductsModel: newProductsModel(conn, c, opts...),
	}
}

func (m *customProductsModel) FindManyByCategoryId(ctx context.Context, categoryId uint64, page, pageSize int) ([]*Products, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf("select %s from %s where `category_id` = ? limit ?,?", productsRows, m.table)
	var resp []*Products
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, categoryId, offset, pageSize)

	return resp, err
}

func (m *customProductsModel) FindOneByName(ctx context.Context, name string) (*Products, error) {
	var resp Products
	query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", productsRows, m.table)
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

func (m *customProductsModel) SearchByKeyword(ctx context.Context, keyword string, page, pageSize int) ([]*Products, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf("select %s from %s where `name` like ? or `description` like ? limit ?,?",
		productsRows, m.table)
	var resp []*Products
	keyword = "%" + keyword + "%"
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, keyword, keyword, offset, pageSize)

	return resp, err
}

func (m *customProductsModel) GeneralSearch(ctx context.Context, page, pageSize int) ([]*Products, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf("select %s from %s limit ?,?",
		productsRows, m.table)
	var resp []*Products
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, offset, pageSize)

	return resp, err
}

func (m *customProductsModel) Count(ctx context.Context, categoryId uint64, keyword string) (int64, error) {
	var count int64
	var query string
	var err error

	if categoryId > 0 && keyword != "" {
		query = fmt.Sprintf("select count(*) from %s where `category_id` = ? and (`name` like ? or `description` like ?)", m.table)
		keyword = "%" + keyword + "%"
		err = m.QueryRowNoCacheCtx(ctx, &count, query, categoryId, keyword, keyword)
	} else if categoryId > 0 {
		query = fmt.Sprintf("select count(*) from %s where `category_id` = ?", m.table)
		err = m.QueryRowNoCacheCtx(ctx, &count, query, categoryId)
	} else if keyword != "" {
		query = fmt.Sprintf("select count(*) from %s where `name` like ? or `description` like ?", m.table)
		keyword = "%" + keyword + "%"
		err = m.QueryRowNoCacheCtx(ctx, &count, query, keyword, keyword)
	} else {
		query = fmt.Sprintf("select count(*) from %s", m.table)
		err = m.QueryRowNoCacheCtx(ctx, &count, query)
	}

	return count, err
}

func (m *customProductsModel) UpdatePartial(ctx context.Context, id uint64, updates map[string]interface{}) error {
    if len(updates) == 0 {
        return nil
    }

    var sets []string
    var args []interface{}
    for k, v := range updates {
        // Special handling for sales increment
        if k == "sales" && strings.Contains(v.(string), "sales + ") {
            sets = append(sets, "`sales` = sales + ?")
            increment := strings.TrimPrefix(v.(string), "sales + ")
            args = append(args, increment)
        } else {
            sets = append(sets, fmt.Sprintf("`%s` = ?", k))
            args = append(args, v)
        }
    }
    args = append(args, id)

    mallProductProductsIdKey := fmt.Sprintf("%s%v", cacheMallProductProductsIdPrefix, id)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(sets, ", "))
        return conn.ExecCtx(ctx, query, args...)
    }, mallProductProductsIdKey)

    return err
}
