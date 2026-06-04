package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func GetOverview(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetOverview(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetHotKnowledgePoints(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := service.GetHotKnowledgePoints(limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetKnowledgeMastery(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetKnowledgeMastery(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetWeakPoints(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := service.GetWeakPoints(userID, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetTrends(c *gin.Context) {
	userID := c.GetUint("user_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	resp, err := service.GetTrends(userID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}
