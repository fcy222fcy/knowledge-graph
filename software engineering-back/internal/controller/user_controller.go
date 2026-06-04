package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetProfile(userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateProfile(userID, req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.ChangePassword(userID, req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := service.ListUsers(page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}
