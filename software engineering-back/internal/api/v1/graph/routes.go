package graph

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册知识图谱模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	graph := protected.Group("/graph")
	{
		graph.GET("", GetGraph)
		graph.POST("/build", BuildGraph)
		graph.GET("/build/latest", GetLatestBuild)
		graph.GET("/build/history", ListBuildHistory)
	}
}
