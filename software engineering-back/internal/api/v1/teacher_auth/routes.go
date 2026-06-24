package teacher_auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册教师认证模块路由（公开接口，无需登录）
func RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/teacher/auth")
	{
		auth.POST("/login", Login)
		auth.POST("/register", Register)
		auth.POST("/refresh", RefreshToken)
	}
}
