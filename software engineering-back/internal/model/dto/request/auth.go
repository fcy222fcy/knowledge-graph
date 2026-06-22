package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名，3-50个字符
	Password string `json:"password" binding:"required,min=6,max=50"` // 密码，6-50个字符
	Email    string `json:"email" binding:"required,email"`           // 邮箱地址
	Nickname string `json:"nickname" binding:"max=50"`               // 昵称，最多50个字符
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	Token string `json:"token" binding:"required"` // 刷新令牌
}
