package request

// CreateKnowledgePointRequest 创建知识点请求
type CreateKnowledgePointRequest struct {
	Name        string `json:"name" binding:"required,max=100"` // 知识点名称，必填
	Description string `json:"description" binding:"max=500"`   // 知识点描述
	DocumentID  uint   `json:"document_id"`                     // 关联的文档ID
	Category    string `json:"category" binding:"max=50"`       // 知识点分类
}

// UpdateKnowledgePointRequest 更新知识点请求
type UpdateKnowledgePointRequest struct {
	Name        string `json:"name" binding:"max=100"` // 知识点名称
	Description string `json:"description" binding:"max=500"` // 知识点描述
	Category    string `json:"category" binding:"max=50"` // 知识点分类
}

// CreateRelationRequest 创建知识点关系请求
type CreateRelationRequest struct {
	SourceID     uint   `json:"source_id" binding:"required"`                                        // 源知识点ID
	TargetID     uint   `json:"target_id" binding:"required"`                                        // 目标知识点ID
	RelationType string `json:"relation_type" binding:"required,oneof=RELATED DEPENDS_ON PART_OF"` // 关系类型：RELATED/DEPENDS_ON/PART_OF
	Description  string `json:"description" binding:"max=500"`                                      // 关系描述
}

// UpdateRelationRequest 更新知识点关系请求
type UpdateRelationRequest struct {
	RelationType string `json:"relation_type" binding:"oneof=RELATED DEPENDS_ON PART_OF"` // 关系类型
	Description  string `json:"description" binding:"max=500"`                             // 关系描述
}
