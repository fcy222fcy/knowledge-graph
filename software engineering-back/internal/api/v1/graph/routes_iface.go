package graph

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册图谱路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, graphService GraphService) {
	ctrl := NewGraphController(graphService)

	graph := protected.Group("/graph")
	{
		graph.GET("", ctrl.GetGraph)
		graph.POST("/build", ctrl.BuildGraph)
		graph.GET("/build/latest", ctrl.GetLatestBuild)
		graph.GET("/build/history", ctrl.ListBuildHistory)
	}
}
