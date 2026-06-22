package quiz

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// SubmitQuiz 提交答题
func SubmitQuiz(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.SubmitQuiz(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, resp)
}

// GetQuizDetail 获取答题记录详情
func GetQuizDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetQuizDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListQuizHistory 获取答题历史记录
func ListQuizHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	list, total, err := service.ListQuizHistory(userID, page, size, uint(knowledgePointID), nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}
