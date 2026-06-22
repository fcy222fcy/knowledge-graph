package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// CreateQuestion 创建题目
func CreateQuestion(c *gin.Context) {
	var req request.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateQuestion(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// GetQuestion 获取题目详情（含答案）
func GetQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetQuestion(uint(id), true)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, resp)
}

// UpdateQuestion 更新题目
func UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateQuestion(uint(id), req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteQuestion 删除题目
func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteQuestion(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListQuestions 获取题目列表
func ListQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	difficulty := c.Query("difficulty")
	list, total, err := service.ListQuestions(page, size, keyword, uint(knowledgePointID), difficulty)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}
