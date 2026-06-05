package question

import (
	"github.com/gin-gonic/gin"
)

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
