package request

type SubmitQuizRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	UserAnswer string `json:"user_answer" binding:"required"`
}
