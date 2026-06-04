package entity

// ─── Ask (Q&A) ────────────────────────────────────────

type AskSession struct {
	BaseModel
	UserID uint   `json:"user_id"`
	Title  string `gorm:"size:200" json:"title"`
}

type AskMessage struct {
	BaseModel
	SessionID  uint    `json:"session_id"`
	Role       string  `gorm:"size:20;not null" json:"role"` // user/assistant
	Content    string  `gorm:"type:text" json:"content"`
	Confidence float64 `json:"confidence"`
}
