package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func SubmitQuiz(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.SubmitQuiz(userID, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetQuizDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetQuizDetail(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListQuizHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	list, total, err := service.ListQuizHistory(userID, page, size, uint(knowledgePointID), nil)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
