package zeroerr

// Common Message Errors
var (
	ErrMessageNotFound     = NewCodeError(106001, "消息不存在")
	ErrMessageCreateFailed = NewCodeError(106002, "创建消息失败")
	ErrMessageUpdateFailed = NewCodeError(106003, "更新消息失败")
	ErrMessageDeleteFailed = NewCodeError(106004, "删除消息失败")
	ErrInvalidMessageType  = NewCodeError(106005, "无效的消息类型")
	ErrInvalidSendChannel  = NewCodeError(106006, "无效的发送渠道")
)

// Template Errors
var (
	ErrTemplateNotFound     = NewCodeError(106101, "消息模板不存在")
	ErrTemplateCreateFailed = NewCodeError(106102, "创建模板失败")
	ErrTemplateUpdateFailed = NewCodeError(106103, "更新模板失败")
	ErrTemplateDeleteFailed = NewCodeError(106104, "删除模板失败")
	ErrDuplicateTemplate    = NewCodeError(106105, "模板代码已存在")
	ErrInvalidTemplate      = NewCodeError(106106, "无效的模板内容")
	ErrTemplateDisabled     = NewCodeError(106107, "模板已禁用")
	ErrInvalidTemplateVars  = NewCodeError(106108, "无效的模板变量")
)

// Message Sending Errors
var (
	ErrSendMessageFailed  = NewCodeError(106201, "发送消息失败")
	ErrBatchSendFailed    = NewCodeError(106202, "批量发送失败")
	ErrChannelSendFailed  = NewCodeError(106203, "渠道发送失败")
	ErrMessageRateLimit   = NewCodeError(106204, "发送频率超限")
	ErrInvalidReceiver    = NewCodeError(106205, "无效的接收者")
	ErrReceiverNotExist   = NewCodeError(106206, "接收者不存在")
	ErrSendQueueFull      = NewCodeError(106207, "发送队列已满")
	ErrRetryLimitExceeded = NewCodeError(106208, "重试次数超限")
)

// Notification Settings Errors
var (
	ErrSettingsNotFound     = NewCodeError(106301, "通知设置不存在")
	ErrSettingsCreateFailed = NewCodeError(106302, "创建通知设置失败")
	ErrSettingsUpdateFailed = NewCodeError(106303, "更新通知设置失败")
	ErrInvalidSettings      = NewCodeError(106304, "无效的通知设置")
	ErrDuplicateSettings    = NewCodeError(106306, "重复的通知设置")
)
