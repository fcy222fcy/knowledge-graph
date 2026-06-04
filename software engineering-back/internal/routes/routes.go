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
				docs.POST("", stubHandler("上传文档"))
				docs.GET("", paginatedStub("文档列表"))
				docs.GET("/:id", stubHandler("文档详情"))
				docs.PUT("/:id", stubHandler("更新文档信息"))
				docs.DELETE("/:id", stubHandler("删除文档"))
				docs.GET("/:id/content", stubHandler("获取文档内容"))
			}

			// Knowledge points
			kp := protected.Group("/knowledge")
			{
				kp.GET("/points", paginatedStub("知识点列表"))
				kp.GET("/points/:id", stubHandler("知识点详情"))
				kp.POST("/points", stubHandler("新增知识点"))
				kp.PUT("/points/:id", stubHandler("更新知识点"))
				kp.DELETE("/points/:id", stubHandler("删除知识点"))
				kp.GET("/relations", paginatedStub("关系列表"))
				kp.POST("/relations", stubHandler("新增关系"))
				kp.PUT("/relations/:id", stubHandler("更新关系"))
				kp.DELETE("/relations/:id", stubHandler("删除关系"))
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
