package response

type KnowledgePointResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type KnowledgeRelationResponse struct {
	ID           uint   `json:"id"`
	SourceID     uint   `json:"source_id"`
	SourceName   string `json:"source_name"`
	TargetID     uint   `json:"target_id"`
	TargetName   string `json:"target_name"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
}
