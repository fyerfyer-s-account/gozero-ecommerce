package zeroerr

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CodeError) Error() string {
	return e.Message
}

func NewCodeError(code int, msg string) error {
	return &CodeError{
		Code:    code,
		Message: msg,
	}
}

// User Common Errors
var (
	ErrUserNotFound      = NewCodeError(100001, "用户不存在")
	ErrDuplicateUsername = NewCodeError(100002, "用户名已存在")
	ErrDuplicatePhone    = NewCodeError(100003, "手机号已存在")
	ErrDuplicateEmail    = NewCodeError(100004, "邮箱已存在")
	ErrInvalidPassword   = NewCodeError(100005, "密码错误")
	ErrInvalidToken      = NewCodeError(100006, "无效的Token")
)

// Authentication Errors
var (
	ErrInvalidAccount  = NewCodeError(100101, "账号不存在")
	ErrAccountDisabled = NewCodeError(100102, "账号已被禁用")
	ErrLoginFailed     = NewCodeError(100103, "登录失败")
	ErrRegisterFailed  = NewCodeError(100104, "注册失败")
)

// User Profile Errors
var (
	ErrUpdateProfile = NewCodeError(100201, "更新用户信息失败")
	ErrInvalidGender = NewCodeError(100202, "无效的性别值")
	ErrInvalidAvatar = NewCodeError(100203, "无效的头像地址")
)

// Address Errors
var (
	ErrAddressNotFound            = NewCodeError(100301, "地址不存在")
	ErrAddressLimit               = NewCodeError(100302, "地址数量超出限制")
	ErrInvalidAddress             = NewCodeError(100303, "无效的地址信息")
	ErrDefaultAddressNotDeletable = NewCodeError(100304, "默认地址不能删除")
)

// Wallet Errors
var (
	ErrInsufficientBalance      = NewCodeError(100401, "余额不足")
	ErrInvalidAmount            = NewCodeError(100402, "无效的金额")
	ErrPayPasswordNotSet        = NewCodeError(100403, "支付密码未设置")
	ErrInvalidPayPassword       = NewCodeError(100404, "支付密码错误")
	ErrWalletDisabled           = NewCodeError(100405, "钱包已被冻结")
	ErrInsufficientFrozenAmount = NewCodeError(100406, "冻结金额不足")
)

// Validation Errors
var (
	ErrInvalidPhone      = NewCodeError(100501, "无效的手机号")
	ErrInvalidEmail      = NewCodeError(100502, "无效的邮箱格式")
	ErrInvalidUsername   = NewCodeError(100503, "无效的用户名")
	ErrPasswordTooWeak   = NewCodeError(100504, "密码强度不足")
	ErrInvalidVerifyCode = NewCodeError(100505, "验证码错误")
)
