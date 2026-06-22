package entity

// ─── Quiz ─────────────────────────────────────────────

// Quiz 答题记录，存储用户的答题结果
type Quiz struct {
	BaseModel
	QuestionID uint   `gorm:"comment:题目ID" json:"question_id"`                         // 题目 ID
	UserID     uint   `gorm:"comment:用户ID" json:"user_id"`                             // 用户 ID
	UserAnswer string `gorm:"size:20;not null;comment:用户提交的答案" json:"user_answer"` // 用户提交的答案
	IsCorrect  bool   `gorm:"comment:是否回答正确" json:"is_correct"`                          // 是否回答正确
}
