package response

// AskResponse 智能问答响应
type AskResponse struct {
	ConversationID         uint        `json:"conversation_id"`         // 会话ID
	QuestionID             uint        `json:"question_id"`             // 问题ID
	Answer                 string      `json:"answer"`                  // LLM 生成的回答
	Confidence             float64     `json:"confidence"`              // 置信度（0-1）
	Sources                []AskSource `json:"sources"`                 // 参考来源
	RelatedKnowledgePoints []KPRef     `json:"related_knowledge_points"` // 相关知识点
	CreatedAt              string      `json:"created_at"`              // 创建时间
}

// AskSource 问答参考来源
type AskSource struct {
	DocumentID    uint   `json:"document_id"`    // 文档ID
	DocumentTitle string `json:"document_title"` // 文档标题
	Content       string `json:"content"`        // 引用的内容片段
}

// KPRef 知识点引用
type KPRef struct {
	ID          uint   `json:"id"`          // 知识点ID
	Name        string `json:"name"`        // 知识点名称
	Description string `json:"description"` // 知识点描述
}

// AskSessionResponse 问答会话响应
type AskSessionResponse struct {
	ConversationID uint   `json:"conversation_id"` // 会话ID
	Title          string `json:"title"`           // 会话标题
	LastQuestion   string `json:"last_question,omitempty"` // 最近问题（可选）
	MessageCount   int    `json:"message_count"`   // 消息数量
	UpdatedAt      string `json:"updated_at"`      // 更新时间
}

// AskMessageResponse 问答消息响应
type AskMessageResponse struct {
	MessageID uint   `json:"message_id"` // 消息ID
	Role      string `json:"role"`       // 角色（user/assistant）
	Content   string `json:"content"`    // 消息内容
	CreatedAt string `json:"created_at"` // 创建时间
}

// AskHistoryItem 问答历史记录项
type AskHistoryItem struct {
	ConversationID uint   `json:"conversation_id"` // 会话ID
	Title          string `json:"title"`           // 会话标题
	LastQuestion   string `json:"last_question"`   // 最近问题
	LastAnswer     string `json:"last_answer"`     // 最近回答
	CreatedAt      string `json:"created_at"`      // 创建时间
	UpdatedAt      string `json:"updated_at"`      // 更新时间
}
