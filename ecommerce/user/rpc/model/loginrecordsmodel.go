package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoginRecordsModel = (*customLoginRecordsModel)(nil)

type (
	// LoginRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginRecordsModel.
	LoginRecordsModel interface {
		loginRecordsModel
		withSession(session sqlx.Session) LoginRecordsModel
		FindByUserId(ctx context.Context, userId uint64, limit int) ([]*LoginRecords, error)
		BatchInsert(ctx context.Context, records []*LoginRecords) error
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
	}

	customLoginRecordsModel struct {
		*defaultLoginRecordsModel
	}
)

// NewLoginRecordsModel returns a model for the database table.
func NewLoginRecordsModel(conn sqlx.SqlConn) LoginRecordsModel {
	return &customLoginRecordsModel{
		defaultLoginRecordsModel: newLoginRecordsModel(conn),
	}
}

func (m *customLoginRecordsModel) FindByUserId(ctx context.Context, userId uint64, limit int) ([]*LoginRecords, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by created_at desc limit ?", loginRecordsRows, m.table)
	var resp []*LoginRecords
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customLoginRecordsModel) BatchInsert(ctx context.Context, records []*LoginRecords) error {
	if len(records) == 0 {
		return nil
	}
	values := make([]string, 0, len(records))
	args := make([]interface{}, 0, len(records)*4)
	for _, record := range records {
		values = append(values, "(?, ?, ?, ?)")
		args = append(args, record.UserId, record.LoginIp, record.LoginLocation, record.DeviceType)
	}
	query := fmt.Sprintf("insert into %s (%s) values %s", m.table, loginRecordsRowsExpectAutoSet, strings.Join(values, ","))
	_, err := m.conn.ExecCtx(ctx, query, args...)
	return err
}

func (m *customLoginRecordsModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customLoginRecordsModel) withSession(session sqlx.Session) LoginRecordsModel {
	return &customLoginRecordsModel{
		defaultLoginRecordsModel: &defaultLoginRecordsModel{
			conn:  sqlx.NewSqlConnFromSession(session),
			table: m.table,
		},
	}
}
