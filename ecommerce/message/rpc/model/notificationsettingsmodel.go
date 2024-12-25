package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ NotificationSettingsModel = (*customNotificationSettingsModel)(nil)

type (
	// NotificationSettingsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNotificationSettingsModel.
	NotificationSettingsModel interface {
		notificationSettingsModel
		withSession(session sqlx.Session) NotificationSettingsModel
	}

	customNotificationSettingsModel struct {
		*defaultNotificationSettingsModel
	}
)

// NewNotificationSettingsModel returns a model for the database table.
func NewNotificationSettingsModel(conn sqlx.SqlConn) NotificationSettingsModel {
	return &customNotificationSettingsModel{
		defaultNotificationSettingsModel: newNotificationSettingsModel(conn),
	}
}

func (m *customNotificationSettingsModel) withSession(session sqlx.Session) NotificationSettingsModel {
	return NewNotificationSettingsModel(sqlx.NewSqlConnFromSession(session))
}
