package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/repository"
	"software_engineering/pkg/response"
)

// ListUsers 获取用户列表（支持搜索和分页）
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	users, total, err := repository.ListUsersAdmin(page, size, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}

	response.Paginated(c, users, total, page, size)
}

// GetUser 获取用户详情
func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.FindUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"status":     user.Status,
		"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at": user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.FindUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Status   *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != nil {
		// 禁止禁用自己
		currentUserID := c.GetUint("user_id")
		if user.ID == currentUserID && *req.Status == 0 {
			response.Error(c, http.StatusBadRequest, "不能禁用自己")
			return
		}
		user.Status = *req.Status
	}

	if err := repository.UpdateUser(user); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, nil)
}

// DeleteUser 删除用户（软删除）
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.FindUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 禁止删除自己
	currentUserID := c.GetUint("user_id")
	if user.ID == currentUserID {
		response.Error(c, http.StatusBadRequest, "不能删除自己")
		return
	}

	if err := repository.DeleteUser(user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 更新用户状态（启用/禁用）
func UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := repository.FindUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 禁止禁用自己
	currentUserID := c.GetUint("user_id")
	if user.ID == currentUserID && req.Status == 0 {
		response.Error(c, http.StatusBadRequest, "不能禁用自己")
		return
	}

	user.Status = req.Status
	if err := repository.UpdateUser(user); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, nil)
}
