package response

type QuizResponse struct {
	QuizID             uint             `json:"quiz_id"`
	QuestionID         uint             `json:"question_id"`
	QuestionTitle      string           `json:"question_title,omitempty"`
	Type               string           `json:"type,omitempty"`
	Difficulty         string           `json:"difficulty,omitempty"`
	Options            []QuestionOption `json:"options,omitempty"`
	UserAnswer         string           `json:"user_answer"`
	CorrectAnswer      string           `json:"correct_answer,omitempty"`
	IsCorrect          bool             `json:"is_correct"`
	Explanation        string           `json:"explanation,omitempty"`
	KnowledgePointID   uint             `json:"knowledge_point_id,omitempty"`
	KnowledgePointName string           `json:"knowledge_point_name,omitempty"`
	CreatedAt          string           `json:"created_at"`
}
