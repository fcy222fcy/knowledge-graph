package document

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutesIface 注册文档路由（依赖注入版本）
func RegisterRoutesIface(protected *gin.RouterGroup, documentService DocumentService) {
	ctrl := NewDocumentController(documentService)

	documents := protected.Group("/documents")
	{
		documents.POST("/upload", ctrl.UploadDocument)
		documents.GET("/:id", ctrl.GetDocument)
		documents.GET("/:id/content", ctrl.GetDocumentContent)
		documents.PUT("/:id", ctrl.UpdateDocument)
		documents.DELETE("/:id", ctrl.DeleteDocument)
		documents.GET("", ctrl.ListDocuments)
	}
}
