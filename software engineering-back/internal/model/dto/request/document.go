package request

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title       string `json:"title" binding:"max=200"`       // 文档标题，最多200个字符
	Description string `json:"description" binding:"max=500"` // 文档描述，最多500个字符
}
