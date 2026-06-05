package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/api/v1/analytics"
	"software_engineering/internal/api/v1/ask"
	"software_engineering/internal/api/v1/auth"
	"software_engineering/internal/api/v1/document"
	"software_engineering/internal/api/v1/graph"
	"software_engineering/internal/api/v1/knowledge"
	"software_engineering/internal/api/v1/question"
	"software_engineering/internal/api/v1/quiz"
	"software_engineering/internal/api/v1/user"
	"software_engineering/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Public routes
		auth.RegisterRoutes(api)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.RequireAuth())
		{
			user.RegisterRoutes(protected)
			document.RegisterRoutes(protected)
			knowledge.RegisterRoutes(protected)
			graph.RegisterRoutes(protected)
			question.RegisterRoutes(protected)
			quiz.RegisterRoutes(protected)
			ask.RegisterRoutes(protected)
			analytics.RegisterRoutes(protected)
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
