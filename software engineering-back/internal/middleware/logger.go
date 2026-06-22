package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件，记录请求方法、路径、状态码和耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		log.Printf("[%s] %s %s %d %v",
			time.Now().Format("2006-01-02 15:04:05"),
			method, path, status, latency)
	}
}
