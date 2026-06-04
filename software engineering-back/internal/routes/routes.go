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
				graph.GET("", stubHandler("获取图谱数据"))
				graph.POST("/build", stubHandler("从文档构建图谱"))
				graph.GET("/build/latest", stubHandler("最近构建结果"))
				graph.GET("/build/history", paginatedStub("构建历史记录"))
			}

			// Questions
			q := protected.Group("/questions")
			{
				q.GET("", paginatedStub("题目列表"))
				q.GET("/:id", stubHandler("题目详情"))
				q.POST("", stubHandler("新增题目"))
				q.PUT("/:id", stubHandler("更新题目"))
				q.DELETE("/:id", stubHandler("删除题目"))
			}

			// Quizzes
			quiz := protected.Group("/quizzes")
			{
				quiz.POST("/submit", stubHandler("提交答题"))
				quiz.GET("/history", paginatedStub("答题历史"))
				quiz.GET("/:id", stubHandler("答题详情"))
			}

			// Ask (Q&A)
			ask := protected.Group("/ask")
			{
				ask.POST("/sessions", stubHandler("新建问答会话"))
				ask.GET("/sessions", paginatedStub("会话列表"))
				ask.GET("/sessions/:id/messages", paginatedStub("会话消息列表"))
				ask.POST("", stubHandler("提问"))
				ask.GET("/history", paginatedStub("问答历史"))
			}

			// Analytics
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/overview", stubHandler("总览统计"))
				analytics.GET("/hot-knowledge-points", stubHandler("热门知识点"))
				analytics.GET("/knowledge-mastery", stubHandler("知识点掌握度"))
				analytics.GET("/weak-points", stubHandler("薄弱知识点"))
				analytics.GET("/trends", stubHandler("趋势数据"))
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
