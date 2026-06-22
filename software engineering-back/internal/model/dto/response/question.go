package response

// QuestionOption 题目选项
type QuestionOption struct {
	Key   string `json:"key"`   // 选项标识（如 A、B、C、D）
	Value string `json:"value"` // 选项内容
}

// QuestionResponse 题目响应
type QuestionResponse struct {
	ID                 uint             `json:"id"`                          // 题目ID
	Title              string           `json:"title"`                       // 题目标题
	Type               string           `json:"type"`                        // 题目类型（single/multiple）
	Difficulty         string           `json:"difficulty"`                  // 难度（easy/medium/hard）
	KnowledgePointID   uint             `json:"knowledge_point_id"`          // 关联的知识点ID
	KnowledgePointName string           `json:"knowledge_point_name,omitempty"` // 知识点名称（可选）
	Options            []QuestionOption `json:"options"`                     // 选项列表
	Answer             string           `json:"answer,omitempty"`            // 正确答案（可选，用于查看答案）
	Explanation        string           `json:"explanation,omitempty"`       // 题目解析（可选）
	CreatedAt          string           `json:"created_at"`                  // 创建时间
}
