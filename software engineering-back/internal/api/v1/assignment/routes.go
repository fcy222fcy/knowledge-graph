package assignment

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/middleware"
)

// RegisterRoutes 注册学生端作业路由
func RegisterRoutes(r *gin.RouterGroup) {
	assignments := r.Group("/assignments")
	assignments.Use(middleware.RequireAuth())
	{
		assignments.GET("", ListAssignments)
		assignments.GET("/:id", GetAssignmentDetail)
		assignments.POST("/:id/submit", SubmitAssignment)
		assignments.GET("/:id/result", GetAssignmentResult)
	}
}

// RegisterAdminRoutes 注册管理端作业路由（挂载到已有的 admin group 上）
func RegisterAdminRoutes(admin *gin.RouterGroup) {
	// admin group 已有 RequireAuth + RequireRole 中间件，不需要重复添加
	assignments := admin.Group("/assignments")
	{
		assignments.GET("", AdminListAssignments)
		assignments.POST("", AdminCreateAssignment)
		assignments.GET("/:id", AdminGetAssignment)
		assignments.PUT("/:id", AdminUpdateAssignment)
		assignments.DELETE("/:id", AdminDeleteAssignment)
		assignments.PUT("/:id/publish", AdminPublishAssignment)
		assignments.PUT("/:id/close", AdminCloseAssignment)
		assignments.GET("/:id/submissions", AdminListSubmissions)
	}
}
