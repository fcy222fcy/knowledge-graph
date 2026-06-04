package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterQuestionRoutes(protected *gin.RouterGroup) {
	q := protected.Group("/questions")
	{
		q.GET("", controller.ListQuestions)
		q.GET("/:id", controller.GetQuestion)
		q.POST("", controller.CreateQuestion)
		q.PUT("/:id", controller.UpdateQuestion)
		q.DELETE("/:id", controller.DeleteQuestion)
	}
}
