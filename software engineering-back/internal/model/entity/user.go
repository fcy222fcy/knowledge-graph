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
	Role     string `gorm:"size:20;default:student;not null;comment:角色 admin/teacher/student" json:"role"`
	Status   int    `gorm:"default:1;comment:状态 1=启用 0=禁用" json:"status"` // 1=启用, 0=禁用
}

// 角色常量
const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
)

// IsAdmin 判断是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsTeacher 判断是否为老师
func (u *User) IsTeacher() bool {
	return u.Role == RoleTeacher
}

// IsStudent 判断是否为学生
func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

// CanAccessAdmin 判断是否可以访问后台
func (u *User) CanAccessAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleTeacher
}
