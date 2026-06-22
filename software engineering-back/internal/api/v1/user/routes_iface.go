package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册用户路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, userService UserService) {
	ctrl := NewUserController(userService)

	users := protected.Group("/users")
	{
		users.GET("/profile", ctrl.GetProfile)
		users.PUT("/profile", ctrl.UpdateProfile)
		users.POST("/password", ctrl.ChangePassword)
		users.GET("/list", ctrl.ListUsers)
	}
}
