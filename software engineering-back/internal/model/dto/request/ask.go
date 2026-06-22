package request

// AskRequest 智能问答请求
type AskRequest struct {
	Question       string `json:"question" binding:"required"` // 用户问题，必填
	ConversationID uint   `json:"conversation_id"`             // 会话ID，用于多轮对话
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Title string `json:"title" binding:"max=200"` // 会话标题，最多200个字符
}
