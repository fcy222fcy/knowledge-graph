package entity

// ─── Knowledge ────────────────────────────────────────

type KnowledgePoint struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `gorm:"size:50" json:"category"`
}

type KnowledgeRelation struct {
	BaseModel
	SourceID     uint   `json:"source_id"`
	TargetID     uint   `json:"target_id"`
	RelationType string `gorm:"size:20;not null" json:"relation_type"` // RELATED/DEPENDS_ON/PART_OF
	Description  string `gorm:"size:500" json:"description"`
}

type KnowledgeBuild struct {
	BaseModel
	DocumentIDs      string `gorm:"size:500" json:"document_ids"` // comma-separated
	CreatedPoints    int    `json:"created_points"`
	CreatedRelations int    `json:"created_relations"`
	ChunkCount       int    `json:"chunk_count"`
	VectorCount      int    `json:"vector_count"`
	Status           string `gorm:"size:20;default:completed" json:"status"`
	Message          string `gorm:"size:500" json:"message"`
}
