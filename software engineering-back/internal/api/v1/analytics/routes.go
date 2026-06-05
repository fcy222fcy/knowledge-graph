package analytics

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup) {
	analytics := protected.Group("/analytics")
	{
		analytics.GET("/overview", GetOverview)
		analytics.GET("/hot-knowledge-points", GetHotKnowledgePoints)
		analytics.GET("/knowledge-mastery", GetKnowledgeMastery)
		analytics.GET("/weak-points", GetWeakPoints)
		analytics.GET("/trends", GetTrends)
	}
}
