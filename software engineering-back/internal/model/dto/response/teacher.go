package response

// TeacherResponse 教师信息响应
type TeacherResponse struct {
	ID        uint   `json:"id"`         // 教师ID
	Username  string `json:"username"`   // 用户名
	Email     string `json:"email"`      // 邮箱
	Nickname  string `json:"nickname"`   // 昵称
	Avatar    string `json:"avatar"`     // 头像URL
	Status    int    `json:"status"`     // 状态（1=启用, 0=禁用）
	CreatedAt string `json:"created_at"` // 创建时间
}

// TeacherLoginResponse 教师登录响应
type TeacherLoginResponse struct {
	Token   string           `json:"token"`   // JWT 访问令牌
	Teacher TeacherResponse  `json:"teacher"` // 教师信息
}
