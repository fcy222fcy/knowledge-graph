package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证模块路由（公开接口，无需登录）
func RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", Register)
		auth.POST("/login", Login)
		auth.POST("/refresh", Refresh)
	}
}
