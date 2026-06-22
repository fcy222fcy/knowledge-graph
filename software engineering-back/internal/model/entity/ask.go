package entity

// ─── Ask (Q&A) ────────────────────────────────────────

// AskSession 问答会话，每次对话会话的容器
type AskSession struct {
	BaseModel
	UserID uint   `gorm:"comment:用户ID" json:"user_id"`
	Title  string `gorm:"size:200;comment:会话标题" json:"title"` // 会话标题，取自第一个问题
}

// AskMessage 问答消息，存储用户提问或助手回答
type AskMessage struct {
	BaseModel
	SessionID  uint    `gorm:"comment:会话ID" json:"session_id"`
	Role       string  `gorm:"size:20;not null;comment:消息角色 user/assistant" json:"role"` // 消息角色：user/assistant
	Content    string  `gorm:"type:text;comment:消息内容" json:"content"`     // 消息内容
	Confidence float64 `gorm:"comment:回答置信度 0-1" json:"confidence"`                   // 回答置信度（0-1）
}
