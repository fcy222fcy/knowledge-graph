package ask

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册问答模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	ask := protected.Group("/ask")
	{
		ask.POST("/sessions", CreateSession)
		ask.GET("/sessions", ListSessions)
		ask.GET("/sessions/:id/messages", ListSessionMessages)
		ask.POST("", AskQuestion)
		ask.GET("/history", ListAskHistory)
	}
}
