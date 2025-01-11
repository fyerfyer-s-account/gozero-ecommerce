package zeroerr

// Common Order Errors
var (
    ErrOrderInvalidParam         = NewCodeError(100001, "无效的参数")
    ErrOrderNotFound        = NewCodeError(100002, "订单不存在")
    ErrOrderNoEmpty         = NewCodeError(100003, "订单号不能为空")
    ErrOrderCreateFailed    = NewCodeError(100004, "创建订单失败")
    ErrOrderUpdateFailed    = NewCodeError(100005, "更新订单失败")
)

// Order Status Errors
var (
    ErrOrderStatusInvalid      = NewCodeError(101001, "无效的订单状态")
    ErrOrderStatusNotAllowed   = NewCodeError(101002, "当前状态不允许此操作")
    ErrOrderAlreadyPaid        = NewCodeError(101003, "订单已支付")
    ErrOrderAlreadyCancelled   = NewCodeError(101004, "订单已取消")
    ErrOrderAlreadyCompleted   = NewCodeError(101005, "订单已完成")
)

// Order Item Errors
var (
    ErrOrderItemNotFound       = NewCodeError(102001, "订单商品不存在")
    ErrOrderItemStockShortage  = NewCodeError(102002, "商品库存不足")
    ErrOrderItemPriceChanged   = NewCodeError(102003, "商品价格已变动")
    ErrOrderItemQuantityLimit  = NewCodeError(102004, "超出商品购买限制")
)

// Order Shipping Errors
var (
    ErrShippingInfoInvalid    = NewCodeError(103001, "配送信息无效")
    ErrShippingNotAvailable   = NewCodeError(103002, "该地区不支持配送")
    ErrShippingStatusInvalid  = NewCodeError(103003, "无效的配送状态")
)

// Order Payment Errors
var (
    ErrPaymentAmountInvalid   = NewCodeError(104001, "支付金额不正确")
    ErrPaymentMethodInvalid   = NewCodeError(104003, "无效的支付方式")
)

// Order Refund Errors
var (
    ErrRefundNotAllowed       = NewCodeError(105002, "订单不允许退款")
    ErrRefundStatusInvalid    = NewCodeError(105003, "无效的退款状态")
)