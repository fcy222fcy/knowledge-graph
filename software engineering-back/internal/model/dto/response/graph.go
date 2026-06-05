package response

type GraphDataResponse struct {
	Nodes   []GraphNode  `json:"nodes"`
	Edges   []GraphEdge  `json:"edges"`
	Summary GraphSummary `json:"summary"`
}

type GraphNode struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `json:"category"`
}

type GraphEdge struct {
	ID           uint   `json:"id"`
	Source       uint   `json:"source"`
	Target       uint   `json:"target"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
}

type GraphSummary struct {
	NodeCount int `json:"node_count"`
	EdgeCount int `json:"edge_count"`
}

type BuildGraphResponse struct {
	BuildID          uint   `json:"build_id"`
	CreatedPoints    int    `json:"created_points"`
	CreatedRelations int    `json:"created_relations"`
	ChunkCount       int    `json:"chunk_count"`
	VectorCount      int    `json:"vector_count"`
	Status           string `json:"status"`
	Message          string `json:"message"`
}

type BuildHistoryResponse struct {
	List      []BuildGraphResponse `json:"list"`
	Total     int64                `json:"total"`
	Page      int                  `json:"page"`
	Size      int                  `json:"size"`
	TotalPage int                  `json:"total_page"`
}
