package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.Register(req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.Login(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, resp)
}

func Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	token, err := service.RefreshToken(req.Token)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	utils.Success(c, gin.H{"token": token})
}
