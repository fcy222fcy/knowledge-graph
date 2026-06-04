package entity

// ─── Quiz ─────────────────────────────────────────────

type Quiz struct {
	BaseModel
	QuestionID uint   `json:"question_id"`
	UserID     uint   `json:"user_id"`
	UserAnswer string `gorm:"size:20;not null" json:"user_answer"`
	IsCorrect  bool   `json:"is_correct"`
}
