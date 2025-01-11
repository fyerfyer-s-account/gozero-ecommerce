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
	ordersFieldNames          = builder.RawFieldNames(&Orders{})
	ordersRows                = strings.Join(ordersFieldNames, ",")
	ordersRowsExpectAutoSet   = strings.Join(stringx.Remove(ordersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	ordersRowsWithPlaceHolder = strings.Join(stringx.Remove(ordersFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallOrderOrdersIdPrefix      = "cache:mallOrder:orders:id:"
	cacheMallOrderOrdersOrderNoPrefix = "cache:mallOrder:orders:orderNo:"
)

type (
	ordersModel interface {
		Insert(ctx context.Context, data *Orders) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*Orders, error)
		FindOneByOrderNo(ctx context.Context, orderNo string) (*Orders, error)
		Update(ctx context.Context, data *Orders) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultOrdersModel struct {
		sqlc.CachedConn
		table string
	}

	Orders struct {
		Id             uint64         `db:"id"`              // è®¢å•ID
		OrderNo        string         `db:"order_no"`        // è®¢å•ç¼–å·
		UserId         uint64         `db:"user_id"`         // ç”¨æˆ·ID
		TotalAmount    float64        `db:"total_amount"`    // è®¢å•æ€»é‡‘é¢
		PayAmount      float64        `db:"pay_amount"`      // åº”ä»˜é‡‘é¢
		FreightAmount  float64        `db:"freight_amount"`  // è¿è´¹
		DiscountAmount float64        `db:"discount_amount"` // ä¼˜æƒ é‡‘é¢
		Status         int64          `db:"status"`          // è®¢å•çŠ¶æ€ 1:å¾…æ”¯ä»˜ 2:å¾…å‘è´§ 3:å¾…æ”¶è´§ 4:å·²å®Œæˆ 5:å·²å–æ¶ˆ 6:å”®åŽä¸­
		Address        string         `db:"address"`         // æ”¶è´§åœ°å€
		Receiver       string         `db:"receiver"`        // æ”¶è´§äºº
		Phone          string         `db:"phone"`           // è”ç³»ç”µè¯
		Remark         sql.NullString `db:"remark"`          // è®¢å•å¤‡æ³¨
		CreatedAt      time.Time      `db:"created_at"`      // åˆ›å»ºæ—¶é—´
		UpdatedAt      time.Time      `db:"updated_at"`      // æ›´æ–°æ—¶é—´
	}
)

func newOrdersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultOrdersModel {
	return &defaultOrdersModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`orders`",
	}
}

func (m *defaultOrdersModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallOrderOrdersIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersIdPrefix, id)
	mallOrderOrdersOrderNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersOrderNoPrefix, data.OrderNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mallOrderOrdersIdKey, mallOrderOrdersOrderNoKey)
	return err
}

func (m *defaultOrdersModel) FindOne(ctx context.Context, id uint64) (*Orders, error) {
	mallOrderOrdersIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersIdPrefix, id)
	var resp Orders
	err := m.QueryRowCtx(ctx, &resp, mallOrderOrdersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", ordersRows, m.table)
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

func (m *defaultOrdersModel) FindOneByOrderNo(ctx context.Context, orderNo string) (*Orders, error) {
	mallOrderOrdersOrderNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersOrderNoPrefix, orderNo)
	var resp Orders
	err := m.QueryRowIndexCtx(ctx, &resp, mallOrderOrdersOrderNoKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `order_no` = ? limit 1", ordersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, orderNo); err != nil {
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

func (m *defaultOrdersModel) Insert(ctx context.Context, data *Orders) (sql.Result, error) {
	mallOrderOrdersIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersIdPrefix, data.Id)
	mallOrderOrdersOrderNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersOrderNoPrefix, data.OrderNo)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, ordersRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.OrderNo, data.UserId, data.TotalAmount, data.PayAmount, data.FreightAmount, data.DiscountAmount, data.Status, data.Address, data.Receiver, data.Phone, data.Remark)
	}, mallOrderOrdersIdKey, mallOrderOrdersOrderNoKey)
	return ret, err
}

func (m *defaultOrdersModel) Update(ctx context.Context, newData *Orders) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	mallOrderOrdersIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersIdPrefix, data.Id)
	mallOrderOrdersOrderNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrdersOrderNoPrefix, data.OrderNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, ordersRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.OrderNo, newData.UserId, newData.TotalAmount, newData.PayAmount, newData.FreightAmount, newData.DiscountAmount, newData.Status, newData.Address, newData.Receiver, newData.Phone, newData.Remark, newData.Id)
	}, mallOrderOrdersIdKey, mallOrderOrdersOrderNoKey)
	return err
}

func (m *defaultOrdersModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallOrderOrdersIdPrefix, primary)
}

func (m *defaultOrdersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", ordersRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultOrdersModel) tableName() string {
	return m.table
}
