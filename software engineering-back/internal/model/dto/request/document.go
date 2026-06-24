package request

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	Title       string `json:"title" binding:"max=200"`       // 文档标题，最多200个字符
	Description string `json:"description" binding:"max=500"` // 文档描述，最多500个字符
}

// ReviewDocumentRequest 审核文档请求
type ReviewDocumentRequest struct {
	Status  string `json:"status" binding:"required,oneof=approved rejected"` // 审核状态：approved 或 rejected
	Comment string `json:"comment"`                                           // 审核意见
}
