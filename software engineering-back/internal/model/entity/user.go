package entity

// ─── User ─────────────────────────────────────────────

type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1=active, 0=disabled
}
