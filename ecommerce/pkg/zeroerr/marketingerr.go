package zeroerr

// Common Marketing Errors
var (
    ErrInvalidMarketingParam = NewCodeError(107001, "无效的营销参数")
    ErrMarketingNotFound     = NewCodeError(107002, "营销活动不存在")
    ErrMarketingExpired      = NewCodeError(107003, "营销活动已过期")
    ErrMarketingNotStarted   = NewCodeError(107004, "营销活动未开始")
    ErrMarketingEnded        = NewCodeError(107005, "营销活动已结束")
    ErrMarketingDisabled     = NewCodeError(107006, "营销活动已禁用")
)

// Coupon Errors
var (
    ErrCouponNotFound        = NewCodeError(107101, "优惠券不存在")
    ErrCouponExpired        = NewCodeError(107102, "优惠券已过期")
    ErrCouponUsed           = NewCodeError(107103, "优惠券已使用")
    ErrCouponReceived       = NewCodeError(107104, "优惠券已领取")
    ErrCouponUnavailable    = NewCodeError(107105, "优惠券已领完")
    ErrCouponNotStarted     = NewCodeError(107106, "优惠券活动未开始")
    ErrCouponEnded          = NewCodeError(107107, "优惠券活动已结束")
    ErrCouponCreateFailed   = NewCodeError(107108, "创建优惠券失败")
    ErrCouponUpdateFailed   = NewCodeError(107109, "更新优惠券失败")
    ErrCouponDeleteFailed   = NewCodeError(107110, "删除优惠券失败")
    ErrExceedCouponLimit    = NewCodeError(107111, "超出领取限制")
    ErrInvalidCouponAmount  = NewCodeError(107112, "无效的优惠券金额")
    ErrMinAmountNotReached  = NewCodeError(107113, "未达到使用门槛")
    ErrCouponNotBelongToUser = NewCodeError(107114, "优惠券不属于该用户")
)

// Promotion Errors
var (
    ErrPromotionNotFound     = NewCodeError(107201, "促销活动不存在")
    ErrPromotionExpired      = NewCodeError(107202, "促销活动已过期")
    ErrPromotionNotStarted   = NewCodeError(107203, "促销活动未开始")
    ErrPromotionEnded        = NewCodeError(107204, "促销活动已结束")
    ErrPromotionDisabled     = NewCodeError(107205, "促销活动已禁用")
    ErrPromotionCreateFailed = NewCodeError(107206, "创建促销活动失败")
    ErrPromotionUpdateFailed = NewCodeError(107207, "更新促销活动失败")
    ErrPromotionDeleteFailed = NewCodeError(107208, "删除促销活动失败")
    ErrInvalidPromotionType  = NewCodeError(107209, "无效的促销类型")
    ErrInvalidPromotionRules = NewCodeError(107210, "无效的促销规则")
    ErrPromotionConflict     = NewCodeError(107211, "促销活动冲突")
    ErrExceedPromotionLimit  = NewCodeError(107212, "超出活动限制")
)

// Points Errors
var (
    ErrPointsNotFound        = NewCodeError(107301, "积分账户不存在")
    ErrInsufficientPoints    = NewCodeError(107302, "积分不足")
    ErrPointsCreateFailed    = NewCodeError(107303, "创建积分账户失败")
    ErrPointsUpdateFailed    = NewCodeError(107304, "更新积分失败")
    ErrPointsDeductFailed    = NewCodeError(107305, "扣减积分失败")
    ErrPointsAddFailed       = NewCodeError(107306, "增加积分失败")
    ErrInvalidPointsAmount   = NewCodeError(107307, "无效的积分数量")
    ErrPointsExpired         = NewCodeError(107308, "积分已过期")
    ErrPointsFreezeFailed    = NewCodeError(107309, "冻结积分失败")
    ErrPointsUnfreezeFailed  = NewCodeError(107310, "解冻积分失败")
    ErrPointsTransferFailed  = NewCodeError(107311, "积分转账失败")
    ErrInvalidPointsOperation = NewCodeError(107312, "无效的积分操作")
    ErrExceedPointsLimit     = NewCodeError(107313, "超出积分限制")
    ErrInvalidPointsSource   = NewCodeError(107314, "无效的积分来源")
)

// Marketing Record Errors
var (
    ErrRecordUpdateFailed    = NewCodeError(107403, "更新记录失败")
    ErrInvalidRecordStatus   = NewCodeError(107406, "无效的记录状态")
)