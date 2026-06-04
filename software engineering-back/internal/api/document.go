package api

import (
	"github.com/gin-gonic/gin"
	"software_engineering/internal/controller"
)

func RegisterDocumentRoutes(protected *gin.RouterGroup) {
	docs := protected.Group("/documents")
	{
		docs.POST("", controller.UploadDocument)
		docs.GET("", controller.ListDocuments)
		docs.GET("/:id", controller.GetDocument)
		docs.PUT("/:id", controller.UpdateDocument)
		docs.DELETE("/:id", controller.DeleteDocument)
		docs.GET("/:id/content", controller.GetDocumentContent)
	}
}
