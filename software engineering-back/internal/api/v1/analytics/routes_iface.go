package analytics

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册分析路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, analyticsService AnalyticsService) {
	ctrl := NewAnalyticsController(analyticsService)

	analytics := protected.Group("/analytics")
	{
		analytics.GET("/overview", ctrl.GetOverview)
		analytics.GET("/hot-knowledge-points", ctrl.GetHotKnowledgePoints)
		analytics.GET("/knowledge-mastery", ctrl.GetKnowledgeMastery)
		analytics.GET("/weak-points", ctrl.GetWeakPoints)
		analytics.GET("/trends", ctrl.GetTrends)
	}
}
