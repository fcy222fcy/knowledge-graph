package graph

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup) {
	graph := protected.Group("/graph")
	{
		graph.GET("", GetGraph)
		graph.POST("/build", BuildGraph)
		graph.GET("/build/latest", GetLatestBuild)
		graph.GET("/build/history", ListBuildHistory)
	}
}
