package knowledge

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册知识点管理模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	kp := protected.Group("/knowledge")
	{
		kp.GET("/points", ListKnowledgePoints)
		kp.GET("/points/:id", GetKnowledgePoint)
		kp.POST("/points", CreateKnowledgePoint)
		kp.PUT("/points/:id", UpdateKnowledgePoint)
		kp.DELETE("/points/:id", DeleteKnowledgePoint)
		kp.GET("/relations", ListRelations)
		kp.POST("/relations", CreateRelation)
		kp.PUT("/relations/:id", UpdateRelation)
		kp.DELETE("/relations/:id", DeleteRelation)
	}
}
