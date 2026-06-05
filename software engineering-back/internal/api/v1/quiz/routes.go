package quiz

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup) {
	quiz := protected.Group("/quizzes")
	{
		quiz.POST("/submit", SubmitQuiz)
		quiz.GET("/history", ListQuizHistory)
		quiz.GET("/:id", GetQuizDetail)
	}
}
