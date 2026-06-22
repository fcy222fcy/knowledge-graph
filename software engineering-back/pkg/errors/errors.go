package errors

import "fmt"

// AppError 应用错误类型，包含错误码和错误消息
type AppError struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误消息
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// New 创建应用错误实例
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Newf 创建格式化的应用错误实例
func Newf(code int, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
