package quiz

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册答题模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	quiz := protected.Group("/quizzes")
	{
		quiz.POST("/submit", SubmitQuiz)
		quiz.GET("/history", ListQuizHistory)
		quiz.GET("/:id", GetQuizDetail)
	}
}
