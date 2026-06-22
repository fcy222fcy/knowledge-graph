package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户管理模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	users := protected.Group("/users")
	{
		users.GET("/profile", GetProfile)
		users.PUT("/profile", UpdateProfile)
		users.POST("/password", ChangePassword)
		users.GET("/list", ListUsers)
	}
}
