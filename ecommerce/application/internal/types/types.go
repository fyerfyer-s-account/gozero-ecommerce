// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Address struct {
	Id            int64  `json:"id"`
	ReceiverName  string `json:"receiverName"`
	ReceiverPhone string `json:"receiverPhone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detailAddress"`
	IsDefault     bool   `json:"isDefault"`
}

type AddressReq struct {
	AddressId     int64  `path:"id"`
	ReceiverName  string `json:"receiverName"`
	ReceiverPhone string `json:"receiverPhone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detailAddress"`
	IsDefault     bool   `json:"isDefault,optional"`
}

type BatchOperateReq struct {
	ItemIds []int64 `json:"itemIds"`
}

type CancelOrderReq struct {
	Id int64 `path:"id"`
}

type CartInfo struct {
	Items         []CartItem `json:"items"`
	TotalPrice    float64    `json:"totalPrice"`
	TotalQuantity int64      `json:"totalQuantity"`
	SelectedPrice float64    `json:"selectedPrice"`
	SelectedCount int64      `json:"selectedCount"`
}

type CartItem struct {
	Id          int64   `json:"id"`
	ProductId   int64   `json:"productId"`
	ProductName string  `json:"productName"`
	SkuId       int64   `json:"skuId"`
	SkuName     string  `json:"skuName"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Selected    bool    `json:"selected"`
	Stock       int32   `json:"stock"`
	CreatedAt   int64   `json:"createdAt"`
}

type CartItemReq struct {
	ProductId int64 `json:"productId"`
	SkuId     int64 `json:"skuId"`
	Quantity  int64 `json:"quantity"`
}

type Category struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ParentId int64  `json:"parentId"`
	Level    int32  `json:"level"`
	Sort     int32  `json:"sort"`
	Icon     string `json:"icon,optional"`
}

type ChangePasswordReq struct {
	UserId      int64  `json:"userId"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ConfirmOrderReq struct {
	Id int64 `path:"id"`
}

type CreateCategoryReq struct {
	Name     string `json:"name"`
	ParentId int64  `json:"parentId"`
	Sort     int32  `json:"sort"`
	Icon     string `json:"icon,optional"`
}

type CreateCategoryResp struct {
	Id int64 `json:"id"`
}

type CreateOrderReq struct {
	AddressId int64  `json:"addressId"`
	Note      string `json:"note,optional"`
}

type CreatePaymentReq struct {
	OrderNo     string  `json:"orderNo"`
	PaymentType int32   `json:"paymentType"`
	Amount      float64 `json:"amount"`
	NotifyUrl   string  `json:"notifyUrl,optional"`
	ReturnUrl   string  `json:"returnUrl,optional"`
}

type CreatePaymentResp struct {
	PaymentNo string `json:"paymentNo"`
	PayUrl    string `json:"payUrl,optional"` // 支付链接或支付参数
	QrCode    string `json:"qrCode,optional"` // 二维码链接
}

type CreateProductReq struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CategoryId  int64             `json:"categoryId"`
	Brand       string            `json:"brand"`
	Images      []string          `json:"images"`
	Price       float64           `json:"price"`
	Attributes  []SkuAttributeReq `json:"skuAttributes"`
}

type CreateProductResp struct {
	Id int64 `json:"id"`
}

type CreateRefundReq struct {
	Id      int64    `path:"id"`
	OrderNo string   `json:"orderNo"`
	Reason  string   `json:"reason"`
	Amount  float64  `json:"amount"`
	Desc    string   `json:"desc,optional"`
	Images  []string `json:"images,optional"`
}

type CreateReviewReq struct {
	ProductId int64    `json:"productId"`
	OrderId   int64    `json:"orderId"`
	Rating    int32    `json:"rating"`
	Content   string   `json:"content"`
	Images    []string `json:"images,optional"`
}

type CreateSkuReq struct {
	ProductId  int64             `path:"productId"`
	SkuCode    string            `json:"skuCode"`
	Price      float64           `json:"price"`
	Stock      int64             `json:"stock"`
	Attributes []SkuAttributeReq `json:"attributes"`
}

type CreateSkuResp struct {
	Id int64 `json:"id"`
}

type DeleteAddressReq struct {
	Id int64 `path:"id"`
}

type DeleteCategoryReq struct {
	Id int64 `json:"id"`
}

type DeleteItemReq struct {
	Id    int64 `path:"id"`
	SkuId int64 `json:"skuId"`
}

type DeleteProductReq struct {
	Id int64 `path:"id"`
}

type DeleteReviewReq struct {
	Id int64 `json:"id"`
}

type GetCategoriesResp struct {
	Categories []Category `json:"categories"`
}

type GetOrderReq struct {
	Id int64 `path:"id"`
}

type GetProductReq struct {
	Id int64 `path:"id"`
}

type GetProductSkusReq struct {
	Id int64 `path:"id"`
}

type GetProfileReq struct {
	Id int64 `path:"id"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutReq struct {
	AccessToken string `json:"accessToken"`
}

type LogoutResp struct {
	Success bool `json:"success"`
}

type Order struct {
	Id             int64       `json:"id"`
	OrderNo        string      `json:"orderNo"`
	UserId         int64       `json:"userId"`
	Status         int32       `json:"status"`         // 1:待支付 2:待发货 3:待收货 4:已完成 5:已取消 6:售后中
	TotalAmount    float64     `json:"totalAmount"`    // 订单总金额
	PayAmount      float64     `json:"payAmount"`      // 实付金额
	FreightAmount  float64     `json:"freightAmount"`  // 运费
	DiscountAmount float64     `json:"discountAmount"` // 优惠金额
	CouponAmount   float64     `json:"couponAmount"`   // 优惠券抵扣
	PointsAmount   float64     `json:"pointsAmount"`   // 积分抵扣
	Items          []OrderItem `json:"items"`
	Address        Address     `json:"address"`
	Payment        Payment     `json:"payment"`
	Shipping       Shipping    `json:"shipping"`
	Note           string      `json:"note"`
	CreatedAt      int64       `json:"createdAt"`
	PayTime        int64       `json:"payTime,optional"`
	ShipTime       int64       `json:"shipTime,optional"`
	ReceiveTime    int64       `json:"receiveTime,optional"`
	FinishTime     int64       `json:"finishTime,optional"`
}

type OrderItem struct {
	Id          int64   `json:"id"`
	ProductId   int64   `json:"productId"`
	ProductName string  `json:"productName"`
	SkuId       int64   `json:"skuId"`
	SkuName     string  `json:"skuName"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int32   `json:"quantity"`
	Amount      float64 `json:"amount"`
}

type OrderListReq struct {
	Status   int32 `form:"status,optional"`
	Page     int32 `form:"page,optional,default=1"`
	PageSize int32 `form:"pageSize,optional,default=20"`
}

type OrderListResp struct {
	List       []Order `json:"list"`
	Total      int64   `json:"total"`
	Page       int32   `json:"page"`
	TotalPages int32   `json:"totalPages"`
}

type OrderProduct struct {
	ProductId int64 `json:"productId"`
	SkuId     int64 `json:"skuId"`
	Quantity  int32 `json:"quantity"`
}

type Payment struct {
	PaymentNo   string  `json:"paymentNo"`
	PaymentType int32   `json:"paymentType"`
	Status      int32   `json:"status"`
	Amount      float64 `json:"amount"`
	PayTime     int64   `json:"payTime,optional"`
}

type PaymentNotifyReq struct {
	PaymentType int32  `json:"paymentType"`
	PaymentNo   string `json:"paymentNo"`
	Data        string `json:"data"` // 原始通知数据
}

type PaymentNotifyResp struct {
	Code    int32  `json:"code"` // 200表示成功
	Message string `json:"message"`
}

type PaymentOrder struct {
	PaymentNo   string  `json:"paymentNo"`
	OrderNo     string  `json:"orderNo"`
	UserId      int64   `json:"userId"`
	Amount      float64 `json:"amount"`
	PaymentType int32   `json:"paymentType"` // 1:微信 2:支付宝 3:余额
	Status      int32   `json:"status"`      // 1:待支付 2:支付中 3:已支付 4:已退款 5:已关闭
	PayTime     int64   `json:"payTime,optional"`
	ExpireTime  int64   `json:"expireTime"`
	CreatedAt   int64   `json:"createdAt"`
}

type PaymentStatusReq struct {
	PaymentNo string `path:"paymentNo"`
}

type PaymentStatusResp struct {
	Status   int32   `json:"status"`
	Amount   float64 `json:"amount"`
	PayTime  int64   `json:"payTime,optional"`
	ErrorMsg string  `json:"errorMsg,optional"`
}

type Product struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Brief        string   `json:"brief"`
	Description  string   `json:"description"`
	CategoryId   int64    `json:"categoryId"`
	CategoryName string   `json:"categoryName"`
	Brand        string   `json:"brand"`
	Images       []string `json:"images"`
	Price        float64  `json:"price"`
	Stock        int32    `json:"stock"`
	Sales        int32    `json:"sales"`
	Rating       float64  `json:"rating"`
	Status       int32    `json:"status"`
	CreatedAt    int64    `json:"createdAt"`
}

type RechargeReq struct {
	Amount      float64 `json:"amount"`
	PaymentType int32   `json:"paymentType"` // 1:微信 2:支付宝
}

type RefundInfo struct {
	Id        int64    `json:"id"`
	OrderId   int64    `json:"orderId"`
	RefundNo  string   `json:"refundNo"`
	Status    int32    `json:"status"` // 0:待处理 1:已同意 2:已拒绝 3:已退款
	Amount    float64  `json:"amount"`
	Reason    string   `json:"reason"`
	Desc      string   `json:"desc"`
	Images    []string `json:"images"`
	CreatedAt int64    `json:"createdAt"`
}

type RefundNotifyReq struct {
	PaymentType int32  `json:"paymentType"`
	RefundNo    string `json:"refundNo"`
	Data        string `json:"data"` // 原始通知数据
}

type RefundNotifyResp struct {
	Code    int32  `json:"code"` // 200表示成功
	Message string `json:"message"`
}

type RefundReq struct {
	PaymentNo string  `json:"paymentNo"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
	NotifyUrl string  `json:"notifyUrl,optional"`
}

type RefundResp struct {
	RefundNo string  `json:"refundNo"`
	Amount   float64 `json:"amount"`
	Status   int32   `json:"status"` // 1:退款中 2:已退款 3:退款失败
}

type RefundStatusReq struct {
	RefundNo string `path:"refundNo"`
}

type RefundStatusResp struct {
	Status     int32   `json:"status"`
	Amount     float64 `json:"amount"`
	Reason     string  `json:"reason"`
	RefundTime int64   `json:"refundTime,optional"`
	ErrorMsg   string  `json:"errorMsg,optional"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone,optional"`
	Email    string `json:"email,optional"`
}

type RegisterResp struct {
	UserId int64 `json:"userId"`
}

type ResetPasswordReq struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
}

type Review struct {
	Id        int64    `json:"id"`
	ProductId int64    `json:"productId"`
	OrderId   int64    `json:"orderId"`
	UserId    int64    `json:"userId"`
	UserName  string   `json:"userName"`
	Rating    int32    `json:"rating"`
	Content   string   `json:"content"`
	Images    []string `json:"images,optional"`
	CreatedAt int64    `json:"createdAt"`
}

type ReviewListReq struct {
	ProductId int64 `form:"productId"`
	Rating    int32 `form:"rating,optional"`
	Page      int32 `form:"page,optional,default=1"`
}

type SearchReq struct {
	Keyword    string   `form:"keyword,optional"`
	CategoryId int64    `form:"categoryId,optional"`
	BrandId    int64    `form:"brandId,optional"`
	PriceMin   float64  `form:"priceMin,optional"`
	PriceMax   float64  `form:"priceMax,optional"`
	Attributes []string `form:"attributes,optional"`
	Sort       string   `form:"sort,optional"`  // price,sales,rating
	Order      string   `form:"order,optional"` // asc,desc
	Page       int32    `form:"page,optional,default=1"`
}

type SearchResp struct {
	List       []Product `json:"list"`
	Total      int64     `json:"total"`
	Page       int32     `json:"page"`
	TotalPages int32     `json:"totalPages"`
}

type SelectedItemsResp struct {
	Items         []CartItem `json:"items"`
	TotalPrice    float64    `json:"totalPrice"`
	TotalQuantity int64      `json:"totalQuantity"`
	ValidStock    bool       `json:"validStock"` // 库存是否足够
}

type Shipping struct {
	ShippingNo  string `json:"shippingNo,optional"`
	Company     string `json:"company,optional"`
	Status      int32  `json:"status"` // 0:待发货 1:已发货 2:已签收
	ShipTime    int64  `json:"shipTime,optional"`
	ReceiveTime int64  `json:"receiveTime,optional"`
}

type Sku struct {
	Id         int64             `json:"id"`
	ProductId  int64             `json:"productId"`
	Name       string            `json:"name"`
	Code       string            `json:"code"`
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Attributes map[string]string `json:"attributes"`
}

type SkuAttributeReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TokenResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type Transaction struct {
	Id        int64   `json:"id"`
	UserId    int64   `json:"userId"`
	OrderId   string  `json:"orderId"`
	Amount    float64 `json:"amount"`
	Type      int64   `json:"type"`
	Status    int64   `json:"status"`
	Remark    string  `json:"remark"`
	CreatedAt int64   `json:"createdAt"`
	UpdatedAt int64   `json:"updatedAt"`
}

type TransactionListReq struct {
	Page     int32 `form:"page,default=1"`
	PageSize int32 `form:"pageSize,default=20"`
	Type     int32 `form:"type,optional"`
}

type TransactionListResp struct {
	List       []Transaction `json:"list"`
	Total      int64         `json:"total"`
	Page       int32         `json:"page"`
	TotalPages int32         `json:"totalPages"`
}

type UpdateCategoryReq struct {
	Id   int64  `path:"id"`
	Name string `json:"name,optional"`
	Sort int32  `json:"sort,optional"`
	Icon string `json:"icon,optional"`
}

type UpdateProductReq struct {
	Id          int64    `path:"id"`
	Name        string   `json:"name,optional"`
	Description string   `json:"description,optional"`
	CategoryId  int64    `json:"categoryId,optional"`
	Brand       string   `json:"brand,optional"`
	Images      []string `json:"images,optional"`
	Price       float64  `json:"price,optional"`
	Status      int64    `json:"status,optional"`
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname,optional"`
	Avatar   string `json:"avatar,optional"`
	Gender   string `json:"gender,optional"`
	Phone    string `json:"phone,optional"`
	Email    string `json:"email,optional"`
}

type UpdateReviewReq struct {
	Id     int64 `path:"id"`
	Status int64 `json:"status"`
}

type UpdateSkuReq struct {
	Id         int64             `path:"id"`
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Attributes map[string]string `json:"attributes"`
}

type UserInfo struct {
	Id          int64   `json:"id"`
	Username    string  `json:"username"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Gender      string  `json:"gender"`
	MemberLevel int32   `json:"memberLevel"`
	Balance     float64 `json:"balance"`
	CreatedAt   int64   `json:"createdAt"`
}

type WalletDetail struct {
	Balance      float64 `json:"balance"`
	Status       int64   `json:"status"`
	FrozenAmount float64 `json:"frozenAmount"`
}
