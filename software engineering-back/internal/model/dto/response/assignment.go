package response

// AssignmentResponse 作业列表响应
type AssignmentResponse struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Chapter      string `json:"chapter,omitempty"`
	Deadline     string `json:"deadline"`
	Status       string `json:"status"`
	TotalScore   int    `json:"total_score"`
	QuestionNum  int    `json:"question_num"`
	SubmitCount  int    `json:"submit_count"`
	TotalCount   int    `json:"total_count"`
	TeacherID    uint   `json:"teacher_id"`
	CreatedAt    string `json:"created_at"`
	IsSubmitted  bool   `json:"is_submitted,omitempty"`  // 学生端：是否已提交
	Score        *int   `json:"score,omitempty"`         // 学生端：得分（已提交时返回）
}

// AssignmentDetailResponse 作业详情响应（学生用，不含答案）
type AssignmentDetailResponse struct {
	ID          uint                         `json:"id"`
	Title       string                       `json:"title"`
	Description string                       `json:"description,omitempty"`
	Chapter     string                       `json:"chapter,omitempty"`
	Deadline    string                       `json:"deadline"`
	TotalScore  int                          `json:"total_score"`
	Status      string                       `json:"status"`
	Questions   []AssignmentQuestionResponse `json:"questions"`
}

// AssignmentQuestionResponse 作业题目响应
type AssignmentQuestionResponse struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	Type        string           `json:"type"`
	Options     []QuestionOption `json:"options"`
	Score       int              `json:"score"`
	SortOrder   int              `json:"sort_order"`
	Answer      string           `json:"answer,omitempty"`      // 仅教师端返回
	Explanation string           `json:"explanation,omitempty"` // 仅教师端返回
}

// AssignmentSubmissionResponse 作业提交记录响应
type AssignmentSubmissionResponse struct {
	ID           uint   `json:"id"`
	AssignmentID uint   `json:"assignment_id"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username,omitempty"`
	Score        int    `json:"score"`
	TotalScore   int    `json:"total_score"`
	Status       string `json:"status"`
	SubmittedAt  string `json:"submitted_at"`
}

// AssignmentResultResponse 学生查看作业结果（含题目和答题情况）
type AssignmentResultResponse struct {
	ID           uint                       `json:"id"`
	AssignmentID uint                       `json:"assignment_id"`
	Title        string                     `json:"title"`
	Score        int                        `json:"score"`
	TotalScore   int                        `json:"total_score"`
	Status       string                     `json:"status"`
	SubmittedAt  string                     `json:"submitted_at"`
	Questions    []AssignmentQuestionResult `json:"questions"`
}

// AssignmentQuestionResult 题目答题结果
type AssignmentQuestionResult struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	Type        string           `json:"type"`
	Options     []QuestionOption `json:"options"`
	Score       int              `json:"score"`
	SortOrder   int              `json:"sort_order"`
	Answer      string           `json:"answer"`       // 标准答案
	Explanation string           `json:"explanation,omitempty"`
	MyAnswer    string           `json:"my_answer"`    // 学生的答案
	IsCorrect   bool             `json:"is_correct"`   // 是否正确
}

// AssignmentSubmitResult 提交作业结果
type AssignmentSubmitResult struct {
	SubmissionID uint   `json:"submission_id"`
	Score        int    `json:"score"`
	TotalScore   int    `json:"total_score"`
	Status       string `json:"status"`
}
