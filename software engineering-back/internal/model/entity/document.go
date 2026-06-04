package entity

// ─── Document ─────────────────────────────────────────

type Document struct {
	BaseModel
	Title       string `gorm:"size:200;not null" json:"title"`
	Description string `gorm:"size:500" json:"description"`
	Filename    string `gorm:"size:200;not null" json:"filename"`
	FilePath    string `gorm:"size:500;not null" json:"-"`
	FileSize    int64  `json:"file_size"`
	FileType    string `gorm:"size:20" json:"file_type"`
	Content     string `gorm:"type:longtext" json:"-"`
	Status      string `gorm:"size:20;default:pending" json:"status"` // pending/processing/completed/failed
}
