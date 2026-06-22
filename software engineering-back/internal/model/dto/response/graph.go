package response

// GraphDataResponse 知识图谱数据响应
type GraphDataResponse struct {
	Nodes   []GraphNode  `json:"nodes"`   // 节点列表
	Edges   []GraphEdge  `json:"edges"`   // 边列表
	Summary GraphSummary `json:"summary"` // 图谱统计摘要
}

// GraphNode 知识图谱节点
type GraphNode struct {
	ID          uint   `json:"id"`          // 节点ID
	Name        string `json:"name"`        // 节点名称
	Description string `json:"description"` // 节点描述
	DocumentID  uint   `json:"document_id"` // 关联的文档ID
	Category    string `json:"category"`    // 节点分类
}

// GraphEdge 知识图谱边（关系）
type GraphEdge struct {
	ID           uint   `json:"id"`             // 边ID
	Source       uint   `json:"source"`         // 源节点ID
	Target       uint   `json:"target"`         // 目标节点ID
	RelationType string `json:"relation_type"`  // 关系类型
	Description  string `json:"description"`    // 关系描述
}

// GraphSummary 知识图谱统计摘要
type GraphSummary struct {
	NodeCount int `json:"node_count"` // 节点数量
	EdgeCount int `json:"edge_count"` // 边数量
}

// BuildGraphResponse 构建知识图谱响应
type BuildGraphResponse struct {
	BuildID          uint   `json:"build_id"`          // 构建记录ID
	CreatedPoints    int    `json:"created_points"`    // 创建的知识点数量
	CreatedRelations int    `json:"created_relations"` // 创建的关系数量
	ChunkCount       int    `json:"chunk_count"`       // 文档分块数量
	VectorCount      int    `json:"vector_count"`      // 向量数量
	Status           string `json:"status"`            // 构建状态
	Message          string `json:"message"`           // 构建结果描述
}

// BuildHistoryResponse 构建历史响应
type BuildHistoryResponse struct {
	List      []BuildGraphResponse `json:"list"`       // 构建记录列表
	Total     int64                `json:"total"`      // 总记录数
	Page      int                  `json:"page"`       // 当前页码
	Size      int                  `json:"size"`       // 每页数量
	TotalPage int                  `json:"total_page"` // 总页数
}
