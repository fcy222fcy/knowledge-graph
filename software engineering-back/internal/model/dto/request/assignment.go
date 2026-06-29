package request

// CreateAssignmentRequest 创建作业请求
type CreateAssignmentRequest struct {
	Title       string                      `json:"title" binding:"required,max=200"`
	Description string                      `json:"description"`
	Chapter     string                      `json:"chapter"`
	Deadline    string                      `json:"deadline" binding:"required"`
	Questions   []CreateAssignmentQuestion  `json:"questions" binding:"required,min=1"`
}

// UpdateAssignmentRequest 更新作业请求
type UpdateAssignmentRequest struct {
	Title       string                     `json:"title" binding:"max=200"`
	Description string                     `json:"description"`
	Chapter     string                     `json:"chapter"`
	Deadline    string                     `json:"deadline"`
	Questions   []UpdateAssignmentQuestion `json:"questions"`
}

// CreateAssignmentQuestion 创建作业题目
type CreateAssignmentQuestion struct {
	Title       string           `json:"title" binding:"required,max=500"`
	Type        string           `json:"type" binding:"required,oneof=single multiple judge"`
	Options     []QuestionOption `json:"options" binding:"required,min=2"`
	Answer      string           `json:"answer" binding:"required"`
	Explanation string           `json:"explanation"`
	Score       int              `json:"score" binding:"min=1"`
	SortOrder   int              `json:"sort_order"`
}

// UpdateAssignmentQuestion 更新作业题目
type UpdateAssignmentQuestion struct {
	ID          uint             `json:"id"` // 0 表示新增
	Title       string           `json:"title" binding:"max=500"`
	Type        string           `json:"type" binding:"oneof=single multiple judge"`
	Options     []QuestionOption `json:"options"`
	Answer      string           `json:"answer"`
	Explanation string           `json:"explanation"`
	Score       int              `json:"score" binding:"min=1"`
	SortOrder   int              `json:"sort_order"`
}

// SubmitAssignmentRequest 提交作业请求
type SubmitAssignmentRequest struct {
	Answers []AssignmentAnswerItem `json:"answers" binding:"required,min=1"`
}

// AssignmentAnswerItem 单题答案
type AssignmentAnswerItem struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}
