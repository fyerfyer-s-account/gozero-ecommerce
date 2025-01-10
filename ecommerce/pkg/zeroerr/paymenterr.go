package zeroerr

// API Common Errors
var (
    ErrInvalidParameter = NewCodeError(100001, "invalid parameter")
    ErrPaymentNoEmpty   = NewCodeError(100002, "payment number cannot be empty")
)

// Common Payment Errors
var (
    ErrInvalidPaymentAmount        = NewCodeError(103001, "无效的支付金额")
    ErrPaymentNotFound      = NewCodeError(103002, "支付订单不存在")
    ErrInvalidPaymentStatus = NewCodeError(103003, "无效的支付状态")
    ErrPaymentExpired       = NewCodeError(103004, "支付订单已过期")
    ErrDuplicatePayment     = NewCodeError(103005, "重复的支付请求")
)

// Payment Channel Errors
var (
    ErrChannelNotFound      = NewCodeError(103101, "支付渠道不存在")
    ErrChannelDisabled      = NewCodeError(103102, "支付渠道已禁用")
    ErrPaymentChannelExists = NewCodeError(103103, "支付渠道已存在")
    ErrInvalidChannelConfig = NewCodeError(103104, "无效的渠道配置")
    ErrChannelUpdateFailed  = NewCodeError(103105, "更新渠道失败")
    ErrUnsupportedChannel   = NewCodeError(103106, "不支持的支付渠道")
)

// Refund Errors
var (
    ErrRefundNotFound      = NewCodeError(103201, "退款订单不存在")
    ErrRefundAmountInvalid = NewCodeError(103202, "退款金额无效")
    ErrRefundExceedAmount  = NewCodeError(103203, "退款金额超出可退金额")
    ErrRefundExpired       = NewCodeError(103204, "退款已过期")
    ErrRefundProcessing    = NewCodeError(103205, "退款正在处理中")
    ErrRefundFailed        = NewCodeError(103206, "退款失败")
)

// Payment Notification Errors
var (
    ErrInvalidNotifyData    = NewCodeError(103301, "无效的通知数据")
    ErrNotifySignInvalid    = NewCodeError(103302, "通知签名验证失败")
    ErrNotifyProcessFailed  = NewCodeError(103303, "通知处理失败")
)