package zeroerr

// Common Cart Errors
var (
    ErrCartNotFound     = NewCodeError(104001, "购物车不存在")
    ErrCartCreateFailed = NewCodeError(104002, "创建购物车失败")
    ErrCartUpdateFailed = NewCodeError(104003, "更新购物车失败")
    ErrCartDeleteFailed = NewCodeError(104004, "清空购物车失败")
)

// Cart Item Errors
var (
    ErrItemNotFound      = NewCodeError(104101, "购物车商品不存在")
    ErrItemAddFailed     = NewCodeError(104102, "添加商品失败")
    ErrItemUpdateFailed  = NewCodeError(104103, "更新商品失败")
    ErrItemDeleteFailed  = NewCodeError(104104, "删除商品失败")
    ErrInvalidQuantity   = NewCodeError(104105, "无效的商品数量")
    ErrExceedMaxItems    = NewCodeError(104106, "超出购物车商品数量限制")
    ErrExceedMaxQuantity = NewCodeError(104107, "超出单个商品数量限制")
    ErrItemDuplicate     = NewCodeError(104108, "商品已存在于购物车")
)

// Cart Statistics Errors
var (
    ErrStatsNotFound     = NewCodeError(104201, "购物车统计数据不存在")
    ErrStatsUpdateFailed = NewCodeError(104202, "更新购物车统计失败")
    ErrCartInvalidAmount = NewCodeError(104203, "无效的金额")
)

// Cart Operation Errors
var (
    ErrSelectFailed     = NewCodeError(104301, "选择商品失败")
    ErrDeselectFailed   = NewCodeError(104302, "取消选择商品失败")
    ErrSelectAllFailed  = NewCodeError(104303, "全选商品失败")
    ErrClearFailed      = NewCodeError(104304, "清空购物车失败")
    ErrMergeCartFailed  = NewCodeError(104305, "合并购物车失败")
    ErrCheckoutFailed   = NewCodeError(104306, "结算购物车失败")
    ErrNoItemsSelected  = NewCodeError(104307, "未选择任何商品")
    ErrItemOutOfStock   = NewCodeError(104308, "商品库存不足")
)