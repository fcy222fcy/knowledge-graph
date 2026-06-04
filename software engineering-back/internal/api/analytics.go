package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterAnalyticsRoutes(protected *gin.RouterGroup) {
	analytics := protected.Group("/analytics")
	{
		analytics.GET("/overview", controller.GetOverview)
		analytics.GET("/hot-knowledge-points", controller.GetHotKnowledgePoints)
		analytics.GET("/knowledge-mastery", controller.GetKnowledgeMastery)
		analytics.GET("/weak-points", controller.GetWeakPoints)
		analytics.GET("/trends", controller.GetTrends)
	}
}
