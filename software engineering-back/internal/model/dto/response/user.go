package response

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint   `json:"id"`         // 用户ID
	Username  string `json:"username"`   // 用户名
	Email     string `json:"email"`      // 邮箱
	Nickname  string `json:"nickname"`   // 昵称
	Avatar    string `json:"avatar"`     // 头像URL
	Role      string `json:"role"`       // 角色: admin, teacher, student
	Status    int    `json:"status"`     // 状态（1=启用, 0=禁用）
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// UserListResponse 用户列表响应（分页）
type UserListResponse struct {
	List      []UserResponse `json:"list"`       // 用户列表
	Total     int64          `json:"total"`      // 总用户数
	Page      int            `json:"page"`       // 当前页码
	Size      int            `json:"size"`       // 每页数量
	TotalPage int            `json:"total_page"` // 总页数
}
