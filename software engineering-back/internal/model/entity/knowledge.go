package entity

// ─── Knowledge ────────────────────────────────────────

// KnowledgePoint 知识点实体，表示知识图谱中的一个节点
type KnowledgePoint struct {
	BaseModel
	Name        string `gorm:"size:100;not null;comment:知识点名称" json:"name"`
	Description string `gorm:"size:500;comment:知识点描述" json:"description"`
	DocumentID  uint   `gorm:"comment:关联的文档ID" json:"document_id"`                        // 关联的文档 ID
	Category    string `gorm:"size:50;comment:知识点分类" json:"category"`            // 知识点分类
}

// KnowledgeRelation 知识点关系实体，表示知识图谱中的边
type KnowledgeRelation struct {
	BaseModel
	SourceID     uint   `gorm:"comment:源知识点ID" json:"source_id"`                          // 源知识点 ID
	TargetID     uint   `gorm:"comment:目标知识点ID" json:"target_id"`                          // 目标知识点 ID
	Type         string `gorm:"size:60;not null;default:'';comment:关系类型标识" json:"type"`
	RelationType string `gorm:"size:20;not null;comment:关系类型 RELATED/DEPENDS_ON/PART_OF" json:"relation_type"` // 关系类型：RELATED/DEPENDS_ON/PART_OF
	Description  string `gorm:"size:500;comment:关系描述" json:"description"`
}

// KnowledgeBuild 知识图谱构建记录，记录每次构建的统计信息
type KnowledgeBuild struct {
	BaseModel
	DocumentIDs      string `gorm:"size:500;comment:逗号分隔的文档ID列表" json:"document_ids"` // 逗号分隔的文档 ID 列表
	CreatedPoints    int    `gorm:"comment:创建的知识点数量" json:"created_points"`               // 创建的知识点数量
	CreatedRelations int    `gorm:"comment:创建的关系数量" json:"created_relations"`            // 创建的关系数量
	ChunkCount       int    `gorm:"comment:文档分块数量" json:"chunk_count"`                  // 文档分块数量
	VectorCount      int    `gorm:"comment:向量数量" json:"vector_count"`                 // 向量数量
	Status           string `gorm:"size:20;default:completed;comment:构建状态" json:"status"` // 构建状态
	Message          string `gorm:"size:500;comment:构建结果描述" json:"message"`      // 构建结果描述
}
