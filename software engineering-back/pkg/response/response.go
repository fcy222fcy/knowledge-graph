package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "software_engineering/pkg/errors"
)

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// SuccessWithMessage 返回带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}

// Error 返回错误响应（保持向后兼容）
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// HandleError 统一错误处理：根据 AppError 自动推断 HTTP 状态码并返回标准错误响应。
// 对于非 AppError 的普通 error，返回 500 内部错误。
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if appErr := apperrors.AsAppError(err); appErr != nil {
		c.JSON(appErr.HTTPStatus(), gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    apperrors.CodeInternalError,
		"message": err.Error(),
		"data":    nil,
	})
}

// Paginated 返回分页响应，包含列表数据和分页元信息
func Paginated(c *gin.Context, list interface{}, total int64, page, size int) {
	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":       list,
			"total":      total,
			"page":       page,
			"size":       size,
			"total_page": totalPage,
		},
	})
}
