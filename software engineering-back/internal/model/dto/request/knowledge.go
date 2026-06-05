package request

type CreateKnowledgePointRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	DocumentID  uint   `json:"document_id"`
	Category    string `json:"category" binding:"max=50"`
}

type UpdateKnowledgePointRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
	Category    string `json:"category" binding:"max=50"`
}

type CreateRelationRequest struct {
	SourceID     uint   `json:"source_id" binding:"required"`
	TargetID     uint   `json:"target_id" binding:"required"`
	RelationType string `json:"relation_type" binding:"required,oneof=RELATED DEPENDS_ON PART_OF"`
	Description  string `json:"description" binding:"max=500"`
}

type UpdateRelationRequest struct {
	RelationType string `json:"relation_type" binding:"oneof=RELATED DEPENDS_ON PART_OF"`
	Description  string `json:"description" binding:"max=500"`
}
