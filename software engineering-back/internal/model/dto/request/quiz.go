package request

// SubmitQuizRequest 提交答题请求
type SubmitQuizRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"` // 题目ID
	UserAnswer string `json:"user_answer" binding:"required"` // 用户答案
}
