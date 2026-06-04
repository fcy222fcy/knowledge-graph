package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterAskRoutes(protected *gin.RouterGroup) {
	ask := protected.Group("/ask")
	{
		ask.POST("/sessions", controller.CreateSession)
		ask.GET("/sessions", controller.ListSessions)
		ask.GET("/sessions/:id/messages", controller.ListSessionMessages)
		ask.POST("", controller.AskQuestion)
		ask.GET("/history", controller.ListAskHistory)
	}
}
