package response

// QuizResponse 答题记录响应
type QuizResponse struct {
	QuizID             uint             `json:"quiz_id"`              // 答题记录ID
	QuestionID         uint             `json:"question_id"`          // 题目ID
	QuestionTitle      string           `json:"question_title,omitempty"` // 题目标题（可选）
	Type               string           `json:"type,omitempty"`       // 题目类型（可选）
	Difficulty         string           `json:"difficulty,omitempty"` // 难度（可选）
	Options            []QuestionOption `json:"options,omitempty"`    // 选项列表（可选）
	UserAnswer         string           `json:"user_answer"`          // 用户答案
	CorrectAnswer      string           `json:"correct_answer,omitempty"` // 正确答案（可选）
	IsCorrect          bool             `json:"is_correct"`           // 是否正确
	Explanation        string           `json:"explanation,omitempty"` // 题目解析（可选）
	KnowledgePointID   uint             `json:"knowledge_point_id,omitempty"` // 知识点ID（可选）
	KnowledgePointName string           `json:"knowledge_point_name,omitempty"` // 知识点名称（可选）
	CreatedAt          string           `json:"created_at"`           // 创建时间
}
