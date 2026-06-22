package question

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册题目管理模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	q := protected.Group("/questions")
	{
		q.GET("", ListQuestions)
		q.GET("/:id", GetQuestion)
		q.POST("", CreateQuestion)
		q.PUT("/:id", UpdateQuestion)
		q.DELETE("/:id", DeleteQuestion)
	}
}
