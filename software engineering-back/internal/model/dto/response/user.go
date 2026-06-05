package response

// User response structs (from user.go)
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserListResponse struct {
	List      []UserResponse `json:"list"`
	Total     int64          `json:"total"`
	Page      int            `json:"page"`
	Size      int            `json:"size"`
	TotalPage int            `json:"total_page"`
}
