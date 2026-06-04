package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func CreateQuestion(c *gin.Context) {
	var req dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateQuestion(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, gin.H{"id": id})
}

func GetQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetQuestion(uint(id), true)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateQuestion(uint(id), req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteQuestion(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	difficulty := c.Query("difficulty")
	list, total, err := service.ListQuestions(page, size, keyword, uint(knowledgePointID), difficulty)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
