package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterQuizRoutes(protected *gin.RouterGroup) {
	quiz := protected.Group("/quizzes")
	{
		quiz.POST("/submit", controller.SubmitQuiz)
		quiz.GET("/history", controller.ListQuizHistory)
		quiz.GET("/:id", controller.GetQuizDetail)
	}
}
