package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// GetProfile 获取当前用户个人信息
func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetProfile(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, resp)
}

// UpdateProfile 更新用户个人信息
func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateProfile(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// ChangePassword 修改用户密码
func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.ChangePassword(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := service.ListUsers(page, size)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, resp)
}
