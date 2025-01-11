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
	orderRefundsFieldNames          = builder.RawFieldNames(&OrderRefunds{})
	orderRefundsRows                = strings.Join(orderRefundsFieldNames, ",")
	orderRefundsRowsExpectAutoSet   = strings.Join(stringx.Remove(orderRefundsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	orderRefundsRowsWithPlaceHolder = strings.Join(stringx.Remove(orderRefundsFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheMallOrderOrderRefundsIdPrefix       = "cache:mallOrder:orderRefunds:id:"
	cacheMallOrderOrderRefundsRefundNoPrefix = "cache:mallOrder:orderRefunds:refundNo:"
)

type (
	orderRefundsModel interface {
		Insert(ctx context.Context, data *OrderRefunds) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*OrderRefunds, error)
		FindOneByRefundNo(ctx context.Context, refundNo string) (*OrderRefunds, error)
		Update(ctx context.Context, data *OrderRefunds) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultOrderRefundsModel struct {
		sqlc.CachedConn
		table string
	}

	OrderRefunds struct {
		Id          uint64         `db:"id"`          // é€€æ¬¾ID
		OrderId     uint64         `db:"order_id"`    // è®¢å•ID
		RefundNo    string         `db:"refund_no"`   // é€€æ¬¾ç¼–å·
		Amount      float64        `db:"amount"`      // é€€æ¬¾é‡‘é¢
		Reason      string         `db:"reason"`      // é€€æ¬¾åŽŸå›
		Status      int64          `db:"status"`      // é€€æ¬¾çŠ¶æ€ 0:å¾…å¤„ç† 1:å·²åŒæ„ 2:å·²æ‹’ç» 3:å·²é€€æ¬¾
		Description sql.NullString `db:"description"` // é—®é¢˜æè¿°
		Images      sql.NullString `db:"images"`      // å›¾ç‰‡å‡­è¯
		Reply       sql.NullString `db:"reply"`       // å¤„ç†å›žå¤
		CreatedAt   time.Time      `db:"created_at"`  // åˆ›å»ºæ—¶é—´
		UpdatedAt   time.Time      `db:"updated_at"`  // æ›´æ–°æ—¶é—´
	}
)

func newOrderRefundsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultOrderRefundsModel {
	return &defaultOrderRefundsModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`order_refunds`",
	}
}

func (m *defaultOrderRefundsModel) Delete(ctx context.Context, id uint64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	mallOrderOrderRefundsIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsIdPrefix, id)
	mallOrderOrderRefundsRefundNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsRefundNoPrefix, data.RefundNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mallOrderOrderRefundsIdKey, mallOrderOrderRefundsRefundNoKey)
	return err
}

func (m *defaultOrderRefundsModel) FindOne(ctx context.Context, id uint64) (*OrderRefunds, error) {
	mallOrderOrderRefundsIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsIdPrefix, id)
	var resp OrderRefunds
	err := m.QueryRowCtx(ctx, &resp, mallOrderOrderRefundsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRefundsRows, m.table)
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

func (m *defaultOrderRefundsModel) FindOneByRefundNo(ctx context.Context, refundNo string) (*OrderRefunds, error) {
	mallOrderOrderRefundsRefundNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsRefundNoPrefix, refundNo)
	var resp OrderRefunds
	err := m.QueryRowIndexCtx(ctx, &resp, mallOrderOrderRefundsRefundNoKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `refund_no` = ? limit 1", orderRefundsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, refundNo); err != nil {
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

func (m *defaultOrderRefundsModel) Insert(ctx context.Context, data *OrderRefunds) (sql.Result, error) {
	mallOrderOrderRefundsIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsIdPrefix, data.Id)
	mallOrderOrderRefundsRefundNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsRefundNoPrefix, data.RefundNo)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, orderRefundsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.OrderId, data.RefundNo, data.Amount, data.Reason, data.Status, data.Description, data.Images, data.Reply)
	}, mallOrderOrderRefundsIdKey, mallOrderOrderRefundsRefundNoKey)
	return ret, err
}

func (m *defaultOrderRefundsModel) Update(ctx context.Context, newData *OrderRefunds) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	mallOrderOrderRefundsIdKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsIdPrefix, data.Id)
	mallOrderOrderRefundsRefundNoKey := fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsRefundNoPrefix, data.RefundNo)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, orderRefundsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.OrderId, newData.RefundNo, newData.Amount, newData.Reason, newData.Status, newData.Description, newData.Images, newData.Reply, newData.Id)
	}, mallOrderOrderRefundsIdKey, mallOrderOrderRefundsRefundNoKey)
	return err
}

func (m *defaultOrderRefundsModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheMallOrderOrderRefundsIdPrefix, primary)
}

func (m *defaultOrderRefundsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRefundsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultOrderRefundsModel) tableName() string {
	return m.table
}
