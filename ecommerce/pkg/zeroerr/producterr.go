package zeroerr

// Common Validation Errors
var (
	ErrInvalidParam = NewCodeError(101001, "无效的参数")
	ErrNotFound     = NewCodeError(101002, "资源不存在")
	ErrNoPermission = NewCodeError(101003, "无权限操作")
)

// Product Errors
var (
	ErrProductNotFound     = NewCodeError(102001, "商品不存在")
	ErrProductDuplicate    = NewCodeError(102002, "商品已存在")
	ErrProductCreateFailed = NewCodeError(102003, "创建商品失败")
	ErrProductUpdateFailed = NewCodeError(102004, "更新商品失败")
	ErrProductDeleteFailed = NewCodeError(102005, "删除商品失败")
	ErrInvalidProductPrice = NewCodeError(102006, "无效的商品价格")
	ErrInvalidProductStock = NewCodeError(102007, "无效的商品库存")
	ErrProductHasOrders    = NewCodeError(102008, "商品存在关联订单")
	ErrProductHasReviews   = NewCodeError(102009, "商品存在评论")
)

// Category Errors
var (
	ErrCategoryNotFound     = NewCodeError(102101, "分类不存在")
	ErrCategoryDuplicate    = NewCodeError(102102, "分类已存在")
	ErrCategoryCreateFailed = NewCodeError(102103, "创建分类失败")
	ErrCategoryHasChildren  = NewCodeError(102104, "分类下存在子分类")
	ErrCategoryHasProducts  = NewCodeError(102105, "分类下存在商品")
	ErrInvalidCategoryLevel = NewCodeError(102106, "无效的分类层级")
	ErrCategoryDeleteFailed = NewCodeError(102107, "删除分类失败")
	ErrCategoryUpdateFailed = NewCodeError(102108, "更新分类失败")
)

// SKU Errors
var (
	ErrSkuNotFound       = NewCodeError(102201, "SKU不存在")
	ErrSkuDuplicate      = NewCodeError(102202, "SKU已存在")
	ErrSkuCreateFailed   = NewCodeError(102203, "创建SKU失败")
	ErrSkuUpdateFailed   = NewCodeError(102204, "更新SKU失败")
	ErrSkuDeleteFailed   = NewCodeError(102205, "删除SKU失败")
	ErrInvalidSkuPrice   = NewCodeError(102206, "无效的SKU价格")
	ErrInvalidSkuStock   = NewCodeError(102207, "无效的SKU库存")
	ErrInvalidAttributes = NewCodeError(102208, "无效的规格属性")
)

// Review Errors
var (
	ErrReviewNotFound      = NewCodeError(102301, "评论不存在")
	ErrReviewCreateFailed  = NewCodeError(102302, "创建评论失败")
	ErrReviewUpdateFailed  = NewCodeError(102303, "更新评论失败")
	ErrReviewDeleteFailed  = NewCodeError(102304, "删除评论失败")
	ErrInvalidRating       = NewCodeError(102305, "无效的评分")
	ErrTooManyImages       = NewCodeError(102306, "图片数量超出限制")
	ErrInvalidReviewLength = NewCodeError(102307, "评论长度不符合要求")
)
