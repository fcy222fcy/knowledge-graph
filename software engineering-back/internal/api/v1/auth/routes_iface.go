package auth

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/service"
)

// RegisterRoutesIface 注册认证路由（依赖注入版本）
func RegisterRoutesIface(api *gin.RouterGroup, authService service.AuthService) {
	ctrl := NewAuthController(authService)

	auth := api.Group("/auth")
	{
		auth.POST("/register", ctrl.Register)
		auth.POST("/login", ctrl.Login)
		auth.POST("/refresh", ctrl.Refresh)
	}
}
