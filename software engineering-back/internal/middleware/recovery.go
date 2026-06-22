package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery panic 恢复中间件，捕获 panic 并返回 500 错误响应
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
