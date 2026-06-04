package dto

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

type UserListResponse struct {
	List      []UserResponse `json:"list"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	Size      int            `json:"size"`
	TotalPage int            `json:"total_page"`
}
