package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterUserRoutes(protected *gin.RouterGroup) {
	users := protected.Group("/users")
	{
		users.GET("/profile", controller.GetProfile)
		users.PUT("/profile", controller.UpdateProfile)
		users.POST("/password", controller.ChangePassword)
		users.GET("/list", controller.ListUsers)
	}
}
