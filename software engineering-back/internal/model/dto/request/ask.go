package request

type AskRequest struct {
	Question       string `json:"question" binding:"required"`
	ConversationID uint   `json:"conversation_id"`
}

type CreateSessionRequest struct {
	Title string `json:"title" binding:"max=200"`
}
