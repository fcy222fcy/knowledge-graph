package admin

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/middleware"
)

// RegisterRoutes 注册后台管理路由，需要 admin 或 teacher 角色
func RegisterRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.RequireRole("admin", "teacher"))
	{
		// 用户管理
		admin.GET("/users", ListUsers)
		admin.GET("/users/:id", GetUser)
		admin.PUT("/users/:id", UpdateUser)
		admin.DELETE("/users/:id", DeleteUser)
		admin.PUT("/users/:id/status", UpdateUserStatus)
		admin.PUT("/users/:id/role", UpdateUserRole)

		// 题目管理
		admin.GET("/questions", ListQuestions)
		admin.POST("/questions", CreateQuestion)
		admin.PUT("/questions/:id", UpdateQuestion)
		admin.DELETE("/questions/:id", DeleteQuestion)

		// 资料管理
		admin.GET("/documents", ListDocuments)
		admin.GET("/documents/:id", GetDocument)
		admin.DELETE("/documents/:id", DeleteDocument)
		admin.PUT("/documents/:id/review", ReviewDocument)

		// 知识点管理
		admin.GET("/knowledge/points", ListKnowledgePoints)
		admin.DELETE("/knowledge/points/:id", DeleteKnowledgePoint)
		admin.GET("/knowledge/relations", ListKnowledgeRelations)
		admin.DELETE("/knowledge/relations/:id", DeleteKnowledgeRelation)

		// 系统统计
		admin.GET("/analytics/overview", GetAnalyticsOverview)
		admin.GET("/analytics/users", GetUserStats)

		// 系统配置
		admin.GET("/system/config", GetSystemConfig)
	}
}
