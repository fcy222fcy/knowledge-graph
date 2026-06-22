package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// Error 返回错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
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
