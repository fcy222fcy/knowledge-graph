package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// ListQuestions 获取题目列表
func ListQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListQuestionsSimple(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// CreateQuestion 创建题目
func CreateQuestion(c *gin.Context) {
	var req request.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	question, err := service.CreateQuestion(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, question)
}

// UpdateQuestion 更新题目
func UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := service.UpdateQuestion(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteQuestion 删除题目
func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := service.DeleteQuestion(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}
