package entity

// ─── User ─────────────────────────────────────────────

// User 用户实体，存储用户账户信息
type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null;comment:用户名" json:"username"`
	Password string `gorm:"size:255;not null;comment:密码哈希" json:"-"`                // 密码字段不序列化到 JSON
	Email    string `gorm:"size:100;uniqueIndex;not null;comment:邮箱地址" json:"email"`
	Nickname string `gorm:"size:50;comment:用户昵称" json:"nickname"`
	Avatar   string `gorm:"size:255;comment:头像URL" json:"avatar"`
	Status   int    `gorm:"default:1;comment:状态 1=启用 0=禁用" json:"status"` // 1=启用, 0=禁用
}
