package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterGraphRoutes(protected *gin.RouterGroup) {
	graph := protected.Group("/graph")
	{
		graph.GET("", controller.GetGraph)
		graph.POST("/build", controller.BuildGraph)
		graph.GET("/build/latest", controller.GetLatestBuild)
		graph.GET("/build/history", controller.ListBuildHistory)
	}
}
