package request

// BuildGraphRequest 构建知识图谱请求
type BuildGraphRequest struct {
	DocumentIDs []uint `json:"document_ids" binding:"required,min=1"` // 文档ID列表，至少选择1个文档
}
