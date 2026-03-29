package errcode

import "net/http"

// Error 统一错误类型（支持自动 HTTP 状态映射）
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Message }

// New 创建错误
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// WithMessage 创建带自定义消息的同 Code 错误
func (e *Error) WithMessage(msg string) *Error {
	return &Error{Code: e.Code, Message: msg}
}

// GetHTTPStatus 根据错误码自动映射 HTTP 状态码
func (e *Error) GetHTTPStatus() int {
	switch {
	case e.Code == 0:
		return http.StatusOK
	case e.Code >= 10001 && e.Code <= 10999:
		return http.StatusUnauthorized
	case e.Code >= 20001 && e.Code <= 20999:
		return http.StatusBadRequest
	case e.Code >= 30001 && e.Code <= 30999:
		return http.StatusForbidden
	case e.Code >= 40001 && e.Code <= 40999:
		return http.StatusBadRequest
	case e.Code >= 50001 && e.Code <= 50999:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// 常用错误码
var (
	// 认证相关 (10xxx → 401)
	ErrInvalidToken      = New(10001, "无效的Token")
	ErrTokenExpired      = New(10002, "Token已过期")
	ErrInvalidCredential = New(10004, "用户名或密码错误")
	ErrAccountDisabled   = New(10005, "账号已被禁用")
	ErrAccountLocked     = New(10008, "账号已被锁定")

	// 资源相关 (20xxx → 400)
	ErrAccountNotFound  = New(20001, "账号不存在")
	ErrAccountExists    = New(20002, "账号已存在")
	ErrPasswordTooShort = New(20005, "密码不符合策略要求")

	// 权限相关 (30xxx → 403)
	ErrPermissionDenied = New(30001, "没有操作权限")

	// 参数相关 (40xxx → 400)
	ErrInvalidParams = New(40001, "请求参数错误")

	// 系统相关 (50xxx → 500)
	ErrInternalServer = New(50001, "服务器内部错误")
	ErrDatabaseError  = New(50002, "数据库操作失败")
)
