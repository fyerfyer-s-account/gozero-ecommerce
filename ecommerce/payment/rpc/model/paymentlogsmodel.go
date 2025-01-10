package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaymentLogsModel = (*customPaymentLogsModel)(nil)

type (
	// PaymentLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaymentLogsModel.
	PaymentLogsModel interface {
		paymentLogsModel
		FindByPaymentNo(ctx context.Context, paymentNo string) ([]*PaymentLogs, error)
		FindByType(ctx context.Context, logType int64) ([]*PaymentLogs, error)
		FindByChannel(ctx context.Context, channel int64) ([]*PaymentLogs, error)
		BatchInsert(ctx context.Context, logs []*PaymentLogs) error
	}

	customPaymentLogsModel struct {
		*defaultPaymentLogsModel
	}
)

// NewPaymentLogsModel returns a model for the database table.
func NewPaymentLogsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaymentLogsModel {
	return &customPaymentLogsModel{
		defaultPaymentLogsModel: newPaymentLogsModel(conn, c, opts...),
	}
}

func (m *customPaymentLogsModel) FindByPaymentNo(ctx context.Context, paymentNo string) ([]*PaymentLogs, error) {
	var logs []*PaymentLogs
	query := fmt.Sprintf("select %s from %s where `payment_no` = ?", paymentLogsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &logs, query, paymentNo)
	return logs, err
}

func (m *customPaymentLogsModel) FindByType(ctx context.Context, logType int64) ([]*PaymentLogs, error) {
	var logs []*PaymentLogs
	query := fmt.Sprintf("select %s from %s where `type` = ?", paymentLogsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &logs, query, logType)
	return logs, err
}

func (m *customPaymentLogsModel) FindByChannel(ctx context.Context, channel int64) ([]*PaymentLogs, error) {
	var logs []*PaymentLogs
	query := fmt.Sprintf("select %s from %s where `channel` = ?", paymentLogsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &logs, query, channel)
	return logs, err
}

func (m *customPaymentLogsModel) BatchInsert(ctx context.Context, logs []*PaymentLogs) error {
	if len(logs) == 0 {
		return nil
	}
	
	query := fmt.Sprintf("insert into %s (%s) values", m.table, paymentLogsRowsExpectAutoSet)
	values := make([]string, 0, len(logs))
	args := make([]interface{}, 0, len(logs)*5)
	
	for _, log := range logs {
		values = append(values, "(?, ?, ?, ?, ?)")
		args = append(args, log.PaymentNo, log.Type, log.Channel, log.RequestData, log.ResponseData)
	}
	
	query = query + " " + strings.Join(values, ",")
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, args...)
	})
	
	return err
}
