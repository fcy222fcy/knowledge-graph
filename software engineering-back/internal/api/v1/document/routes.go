package document

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文档管理模块路由
func RegisterRoutes(protected *gin.RouterGroup) {
	docs := protected.Group("/documents")
	{
		docs.POST("", UploadDocument)
		docs.GET("", ListDocuments)
		docs.GET("/:id", GetDocument)
		docs.PUT("/:id", UpdateDocument)
		docs.DELETE("/:id", DeleteDocument)
		docs.GET("/:id/content", GetDocumentContent)
	}
}
