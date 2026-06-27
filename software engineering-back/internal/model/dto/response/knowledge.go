package response

import "software_engineering/internal/model/entity"

// KnowledgePointResponse 知识点响应
type KnowledgePointResponse struct {
	ID          uint                    `json:"id"`           // 知识点ID
	Name        string                  `json:"name"`         // 知识点名称
	Description string                  `json:"description"`  // 知识点描述
	DocumentID  uint                    `json:"document_id"`  // 关联的文档ID
	Category    string                  `json:"category"`     // 知识点分类
	Sources     []entity.DocumentSource `json:"sources"`      // 所有来源文档
	CreatedAt   string                  `json:"created_at"`   // 创建时间
	UpdatedAt   string                  `json:"updated_at"`   // 更新时间
}

// KnowledgeRelationResponse 知识点关系响应
type KnowledgeRelationResponse struct {
	ID           uint   `json:"id"`             // 关系ID
	SourceID     uint   `json:"source_id"`      // 源知识点ID
	SourceName   string `json:"source_name"`    // 源知识点名称
	TargetID     uint   `json:"target_id"`      // 目标知识点ID
	TargetName   string `json:"target_name"`    // 目标知识点名称
	RelationType string `json:"relation_type"`  // 关系类型
	Description  string `json:"description"`    // 关系描述
	CreatedAt    string `json:"created_at"`     // 创建时间
}
