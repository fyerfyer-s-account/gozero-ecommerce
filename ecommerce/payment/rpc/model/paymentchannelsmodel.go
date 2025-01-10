package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentChannelsModel = (*customPaymentChannelsModel)(nil)

type (
	// PaymentChannelsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentChannelsModel.
	PaymentChannelsModel interface {
		paymentChannelsModel
		FindAll(ctx context.Context) ([]*PaymentChannels, error)
		FindManyByStatus(ctx context.Context, status int64) ([]*PaymentChannels, error)
		UpdateFields(ctx context.Context, id uint64, updates map[string]interface{}) error
		FindOneByChannelAndStatus(ctx context.Context, channel, status int64) (*PaymentChannels, error)
	}

	customPaymentChannelsModel struct {
		*defaultPaymentChannelsModel
	}
)

// NewPaymentChannelsModel returns a model for the database table.
func NewPaymentChannelsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentChannelsModel {
	return &customPaymentChannelsModel{
		defaultPaymentChannelsModel: newPaymentChannelsModel(conn, c, opts...),
	}
}

func (m *customPaymentChannelsModel) FindAll(ctx context.Context) ([]*PaymentChannels, error) {
	var channels []*PaymentChannels
	query := fmt.Sprintf("select %s from %s", paymentChannelsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &channels, query)
	return channels, err
}

func (m *customPaymentChannelsModel) FindManyByStatus(ctx context.Context, status int64) ([]*PaymentChannels, error) {
	var channels []*PaymentChannels
	query := fmt.Sprintf("select %s from %s where `status` = ?", paymentChannelsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &channels, query, status)
	return channels, err
}

func (m *customPaymentChannelsModel) UpdateFields(ctx context.Context, id uint64, updates map[string]interface{}) error {
	oldData, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	var sets []string
	var args []interface{}
	for k, v := range updates {
		sets = append(sets, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	args = append(args, id)

	mallPaymentPaymentChannelsChannelKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentChannelsChannelPrefix, oldData.Channel)
	mallPaymentPaymentChannelsIdKey := fmt.Sprintf("%s%v", cacheMallPaymentPaymentChannelsIdPrefix, id)

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(sets, ", "))
		return conn.ExecCtx(ctx, query, args...)
	}, mallPaymentPaymentChannelsChannelKey, mallPaymentPaymentChannelsIdKey)

	return err
}

func (m *customPaymentChannelsModel) FindOneByChannelAndStatus(ctx context.Context, channel, status int64) (*PaymentChannels, error) {
	var resp PaymentChannels
	query := fmt.Sprintf("select %s from %s where `channel` = ? and `status` = ? limit 1", paymentChannelsRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, channel, status)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
