package knowledge

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册知识点路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, knowledgeService KnowledgeService) {
	ctrl := NewKnowledgeController(knowledgeService)

	knowledge := protected.Group("/knowledge")
	{
		knowledge.POST("/points", ctrl.CreateKnowledgePoint)
		knowledge.GET("/points/:id", ctrl.GetKnowledgePoint)
		knowledge.PUT("/points/:id", ctrl.UpdateKnowledgePoint)
		knowledge.DELETE("/points/:id", ctrl.DeleteKnowledgePoint)
		knowledge.GET("/points", ctrl.ListKnowledgePoints)
		knowledge.POST("/relations", ctrl.CreateRelation)
		knowledge.PUT("/relations/:id", ctrl.UpdateRelation)
		knowledge.DELETE("/relations/:id", ctrl.DeleteRelation)
		knowledge.GET("/relations", ctrl.ListRelations)
	}
}
