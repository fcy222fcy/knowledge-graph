package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
	"software_engineering/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.Register)
			auth.POST("/login", controller.Login)
			auth.POST("/refresh", controller.Refresh)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.RequireAuth())
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/profile", controller.GetProfile)
				users.PUT("/profile", controller.UpdateProfile)
				users.POST("/password", controller.ChangePassword)
				users.GET("/list", controller.ListUsers)
			}

			// Documents
			docs := protected.Group("/documents")
			{
				docs.POST("", controller.UploadDocument)
				docs.GET("", controller.ListDocuments)
				docs.GET("/:id", controller.GetDocument)
				docs.PUT("/:id", controller.UpdateDocument)
				docs.DELETE("/:id", controller.DeleteDocument)
				docs.GET("/:id/content", controller.GetDocumentContent)
			}

			// Knowledge points
			kp := protected.Group("/knowledge")
			{
				kp.GET("/points", controller.ListKnowledgePoints)
				kp.GET("/points/:id", controller.GetKnowledgePoint)
				kp.POST("/points", controller.CreateKnowledgePoint)
				kp.PUT("/points/:id", controller.UpdateKnowledgePoint)
				kp.DELETE("/points/:id", controller.DeleteKnowledgePoint)
				kp.GET("/relations", controller.ListRelations)
				kp.POST("/relations", controller.CreateRelation)
				kp.PUT("/relations/:id", controller.UpdateRelation)
				kp.DELETE("/relations/:id", controller.DeleteRelation)
			}

			// Graph
			graph := protected.Group("/graph")
			{
				graph.GET("", controller.GetGraph)
				graph.POST("/build", controller.BuildGraph)
				graph.GET("/build/latest", controller.GetLatestBuild)
				graph.GET("/build/history", controller.ListBuildHistory)
			}

			// Questions
			q := protected.Group("/questions")
			{
				q.GET("", controller.ListQuestions)
				q.GET("/:id", controller.GetQuestion)
				q.POST("", controller.CreateQuestion)
				q.PUT("/:id", controller.UpdateQuestion)
				q.DELETE("/:id", controller.DeleteQuestion)
			}

			// Quizzes
			quiz := protected.Group("/quizzes")
			{
				quiz.POST("/submit", controller.SubmitQuiz)
				quiz.GET("/history", controller.ListQuizHistory)
				quiz.GET("/:id", controller.GetQuizDetail)
			}

			// Ask (Q&A)
			ask := protected.Group("/ask")
			{
				ask.POST("/sessions", controller.CreateSession)
				ask.GET("/sessions", controller.ListSessions)
				ask.GET("/sessions/:id/messages", controller.ListSessionMessages)
				ask.POST("", controller.AskQuestion)
				ask.GET("/history", controller.ListAskHistory)
			}

			// Analytics
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/overview", controller.GetOverview)
				analytics.GET("/hot-knowledge-points", controller.GetHotKnowledgePoints)
				analytics.GET("/knowledge-mastery", controller.GetKnowledgeMastery)
				analytics.GET("/weak-points", controller.GetWeakPoints)
				analytics.GET("/trends", controller.GetTrends)
			}
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

func stubHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": name + " - 待实现",
			"data":    nil,
		})
	}
}

func paginatedStub(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": name + " - 待实现",
			"data": gin.H{
				"list":       []interface{}{},
				"total":      0,
				"page":       page,
				"size":       size,
				"total_page": 0,
			},
		})
	}
}
