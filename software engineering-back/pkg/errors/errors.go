package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用错误类型，包含错误码和错误消息
type AppError struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误消息
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// HTTPStatus 根据错误码返回对应的 HTTP 状态码
func (e *AppError) HTTPStatus() int {
	switch {
	case e.Code == CodeSuccess:
		return http.StatusOK
	case e.Code == CodeBadRequest || e.Code == CodeValidationFailed:
		return http.StatusBadRequest
	case e.Code == CodeUnauthorized || e.Code == CodeInvalidToken || e.Code == CodeInvalidPassword || e.Code == CodeUserDisabled:
		return http.StatusUnauthorized
	case e.Code == CodeForbidden || e.Code == CodeDocumentAccessDenied:
		return http.StatusForbidden
	case e.Code == CodeNotFound || e.Code == CodeUserNotFound || e.Code == CodeDocumentNotFound ||
		e.Code == CodeKnowledgePointNotFound || e.Code == CodeRelationNotFound ||
		e.Code == CodeQuestionNotFound || e.Code == CodeQuizNotFound || e.Code == CodeAskSessionNotFound:
		return http.StatusNotFound
	case e.Code >= 7000:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
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

// IsAppError 检查 error 是否为 AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError 将 error 转换为 AppError，非 AppError 时返回内部错误
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return New(CodeInternalError, err.Error())
}
