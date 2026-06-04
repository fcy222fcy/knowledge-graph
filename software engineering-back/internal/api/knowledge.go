package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterKnowledgeRoutes(protected *gin.RouterGroup) {
	kp := protected.Group("/knowledge")
	{
		kp.GET("/points", controller.ListKnowledgePoints)
		kp.GET("/points/:id", controller.GetKnowledgePoint)
		kp.POST("/points", controller.CreateKnowledgePoint)
		kp.PUT("/points/:id", controller.UpdateKnowledgePoint)
		kp.DELETE("/points/:id", controller.DeleteKnowledgePoint)
		kp.GET("/relations", controller.ListRelations)
		kp.POST("/relations", controller.CreateRelation)
		kp.PUT("/relations/:id", controller.UpdateRelation)
		kp.DELETE("/relations/:id", controller.DeleteRelation)
	}
}
