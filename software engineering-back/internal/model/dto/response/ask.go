package response

type AskResponse struct {
	ConversationID         uint        `json:"conversation_id"`
	QuestionID             uint        `json:"question_id"`
	Answer                 string      `json:"answer"`
	Confidence             float64     `json:"confidence"`
	Sources                []AskSource `json:"sources"`
	RelatedKnowledgePoints []KPRef     `json:"related_knowledge_points"`
	CreatedAt              string      `json:"created_at"`
}

type AskSource struct {
	DocumentID    uint   `json:"document_id"`
	DocumentTitle string `json:"document_title"`
	Content       string `json:"content"`
}

type KPRef struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AskSessionResponse struct {
	ConversationID uint   `json:"conversation_id"`
	Title          string `json:"title"`
	LastQuestion   string `json:"last_question,omitempty"`
	MessageCount   int    `json:"message_count"`
	UpdatedAt      string `json:"updated_at"`
}

type AskMessageResponse struct {
	MessageID uint   `json:"message_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
