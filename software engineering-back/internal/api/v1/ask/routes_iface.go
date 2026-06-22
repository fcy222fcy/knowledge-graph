package ask

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册问答路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, askService AskService) {
	ctrl := NewAskController(askService)

	ask := protected.Group("/ask")
	{
		ask.POST("/sessions", ctrl.CreateSession)
		ask.GET("/sessions", ctrl.ListSessions)
		ask.GET("/sessions/:id/messages", ctrl.ListSessionMessages)
		ask.POST("/question", ctrl.AskQuestion)
		ask.GET("/history", ctrl.ListAskHistory)
	}
}
