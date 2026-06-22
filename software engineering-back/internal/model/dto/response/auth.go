package response

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"` // JWT 访问令牌
	User  UserResponse `json:"user"`  // 用户信息
}
