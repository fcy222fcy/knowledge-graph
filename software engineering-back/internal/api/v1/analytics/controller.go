package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// GetOverview 获取学习概览数据
func GetOverview(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetOverview(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetHotKnowledgePoints 获取热门知识点列表
func GetHotKnowledgePoints(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := service.GetHotKnowledgePoints(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetKnowledgeMastery 获取知识点掌握程度
func GetKnowledgeMastery(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetKnowledgeMastery(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetWeakPoints 获取薄弱知识点
func GetWeakPoints(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := service.GetWeakPoints(userID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetTrends 获取学习趋势数据
func GetTrends(c *gin.Context) {
	userID := c.GetUint("user_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	resp, err := service.GetTrends(userID, days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}
