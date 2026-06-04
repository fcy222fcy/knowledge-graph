package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// 注册各模块路由
		RegisterAuthRoutes(api)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.RequireAuth())
		{
			RegisterUserRoutes(protected)
			RegisterDocumentRoutes(protected)
			RegisterKnowledgeRoutes(protected)
			RegisterGraphRoutes(protected)
			RegisterQuestionRoutes(protected)
			RegisterQuizRoutes(protected)
			RegisterAskRoutes(protected)
			RegisterAnalyticsRoutes(protected)
		}
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"status":  "ok",
			"service": "software-engineering-backend",
		},
	})
}
