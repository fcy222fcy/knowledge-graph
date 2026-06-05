package request

type CreateQuestionRequest struct {
	Title            string           `json:"title" binding:"required,max=500"`
	Type             string           `json:"type" binding:"required,oneof=single multiple"`
	Difficulty       string           `json:"difficulty" binding:"required,oneof=easy medium hard"`
	KnowledgePointID uint             `json:"knowledge_point_id" binding:"required"`
	Options          []QuestionOption `json:"options" binding:"required,min=2"`
	Answer           string           `json:"answer" binding:"required"`
	Explanation      string           `json:"explanation"`
}

type UpdateQuestionRequest struct {
	Title       string           `json:"title" binding:"max=500"`
	Type        string           `json:"type" binding:"oneof=single multiple"`
	Difficulty  string           `json:"difficulty" binding:"oneof=easy medium hard"`
	Options     []QuestionOption `json:"options"`
	Answer      string           `json:"answer"`
	Explanation string           `json:"explanation"`
}

type QuestionOption struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
