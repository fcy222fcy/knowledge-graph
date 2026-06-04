package entity

// ─── Question ─────────────────────────────────────────

type Question struct {
	BaseModel
	Title            string `gorm:"size:500;not null" json:"title"`
	Type             string `gorm:"size:20;not null" json:"type"`       // single/multiple
	Difficulty       string `gorm:"size:20;not null" json:"difficulty"` // easy/medium/hard
	KnowledgePointID uint   `json:"knowledge_point_id"`
	Options          string `gorm:"type:text" json:"-"` // JSON array stored as text
	Answer           string `gorm:"size:20;not null" json:"answer"`
	Explanation      string `gorm:"type:text" json:"explanation"`
}
