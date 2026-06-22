package request

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"` // 昵称，最多50个字符
	Avatar   string `json:"avatar" binding:"max=255"` // 头像URL，最多255个字符
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`           // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"` // 新密码，6-50个字符
}
