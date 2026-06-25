package entity

// Teacher 教师实体（相当于管理员）
type Teacher struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null;comment:用户名" json:"username"`
	Password string `gorm:"size:255;not null;comment:密码哈希" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex;not null;comment:邮箱地址" json:"email"`
	Nickname string `gorm:"size:50;comment:教师昵称" json:"nickname"`
	Avatar   string `gorm:"size:255;comment:头像URL" json:"avatar"`
	Status   int    `gorm:"default:1;comment:状态 1=启用 0=禁用" json:"status"`
}
