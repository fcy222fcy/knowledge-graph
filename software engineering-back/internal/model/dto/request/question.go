package request

// CreateQuestionRequest 创建题目请求
type CreateQuestionRequest struct {
	Title            string           `json:"title" binding:"required,max=500"`    // 题目标题，必填
	Type             string           `json:"type" binding:"required,oneof=single multiple"` // 题目类型：single(单选)/multiple(多选)
	Difficulty       string           `json:"difficulty" binding:"required,oneof=easy medium hard"` // 难度：easy/medium/hard
	KnowledgePointID uint             `json:"knowledge_point_id" binding:"required"` // 关联的知识点ID
	Options          []QuestionOption `json:"options" binding:"required,min=2"`     // 选项列表，至少2个选项
	Answer           string           `json:"answer" binding:"required"`           // 正确答案
	Explanation      string           `json:"explanation"`                          // 题目解析
}

// UpdateQuestionRequest 更新题目请求
type UpdateQuestionRequest struct {
	Title       string           `json:"title" binding:"max=500"`                  // 题目标题
	Type        string           `json:"type" binding:"oneof=single multiple"`     // 题目类型
	Difficulty  string           `json:"difficulty" binding:"oneof=easy medium hard"` // 难度
	Options     []QuestionOption `json:"options"`                                  // 选项列表
	Answer      string           `json:"answer"`                                   // 正确答案
	Explanation string           `json:"explanation"`                              // 题目解析
}

// QuestionOption 题目选项
type QuestionOption struct {
	Key   string `json:"key"`   // 选项标识（如 A、B、C、D）
	Value string `json:"value"` // 选项内容
}
