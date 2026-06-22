package question

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册题目路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, questionService QuestionService) {
	ctrl := NewQuestionController(questionService)

	questions := protected.Group("/questions")
	{
		questions.POST("", ctrl.CreateQuestion)
		questions.GET("/:id", ctrl.GetQuestion)
		questions.PUT("/:id", ctrl.UpdateQuestion)
		questions.DELETE("/:id", ctrl.DeleteQuestion)
		questions.GET("", ctrl.ListQuestions)
	}
}
