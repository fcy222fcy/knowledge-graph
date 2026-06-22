package quiz

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册答题路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, quizService QuizService) {
	ctrl := NewQuizController(quizService)

	quizzes := protected.Group("/quizzes")
	{
		quizzes.POST("/submit", ctrl.SubmitQuiz)
		quizzes.GET("/:id", ctrl.GetQuizDetail)
		quizzes.GET("/history", ctrl.ListQuizHistory)
	}
}
