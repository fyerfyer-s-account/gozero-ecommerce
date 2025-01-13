package model

import (
    "context"
    "database/sql"
    "fmt"
    "strings"

    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NotificationSettingsModel = (*customNotificationSettingsModel)(nil)

type (
    NotificationSettingsModel interface {
        notificationSettingsModel
        Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
        FindByUserId(ctx context.Context, userId uint64) ([]*NotificationSettings, error)
        FindByUserIdAndType(ctx context.Context, userId uint64, tp int64) ([]*NotificationSettings, error)
        BatchUpsert(ctx context.Context, data []*NotificationSettings) error
        UpdateSettings(ctx context.Context, userId uint64, tp, channel int64, isEnabled bool) error
    }

    customNotificationSettingsModel struct {
        *defaultNotificationSettingsModel
    }
)

func NewNotificationSettingsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) NotificationSettingsModel {
    return &customNotificationSettingsModel{
        defaultNotificationSettingsModel: newNotificationSettingsModel(conn, c, opts...),
    }
}

func (m *customNotificationSettingsModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
    return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
        return fn(ctx, session)
    })
}

func (m *customNotificationSettingsModel) FindByUserId(ctx context.Context, userId uint64) ([]*NotificationSettings, error) {
    var settings []*NotificationSettings
    query := fmt.Sprintf("select %s from %s where `user_id` = ?", notificationSettingsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &settings, query, userId)
    return settings, err
}

func (m *customNotificationSettingsModel) FindByUserIdAndType(ctx context.Context, userId uint64, tp int64) ([]*NotificationSettings, error) {
    var settings []*NotificationSettings
    query := fmt.Sprintf("select %s from %s where `user_id` = ? and `type` = ?", notificationSettingsRows, m.table)
    err := m.QueryRowsNoCacheCtx(ctx, &settings, query, userId, tp)
    return settings, err
}

func (m *customNotificationSettingsModel) BatchUpsert(ctx context.Context, data []*NotificationSettings) error {
    if len(data) == 0 {
        return nil
    }

    values := make([]string, 0, len(data))
    args := make([]interface{}, 0, len(data)*4)
    for _, setting := range data {
        values = append(values, "(?, ?, ?, ?)")
        args = append(args, setting.UserId, setting.Type, setting.Channel, setting.IsEnabled)
    }

    query := fmt.Sprintf("insert into %s (%s) values %s on duplicate key update "+
        "`is_enabled` = values(`is_enabled`)", 
        m.table, notificationSettingsRowsExpectAutoSet, strings.Join(values, ","))

    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
        return conn.ExecCtx(ctx, query, args...)
    })
    return err
}

func (m *customNotificationSettingsModel) UpdateSettings(ctx context.Context, userId uint64, tp, channel int64, isEnabled bool) error {
    enabled := int64(0)
    if isEnabled {
        enabled = 1
    }

    setting, err := m.FindOneByUserIdTypeChannel(ctx, userId, tp, channel)
    if err != nil && err != ErrNotFound {
        return err
    }

    if err == ErrNotFound {
        // Insert new setting
        _, err = m.Insert(ctx, &NotificationSettings{
            UserId:    userId,
            Type:      tp,
            Channel:   channel,
            IsEnabled: enabled,
        })
        return err
    }

    // Update existing setting
    setting.IsEnabled = enabled
    return m.Update(ctx, setting)
}