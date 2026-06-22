package ask

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// CreateSession 创建新的问答会话
func CreateSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.CreateSession(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListSessions 获取用户的问答会话列表
func ListSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := service.ListSessions(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// ListSessionMessages 获取指定会话的消息记录
func ListSessionMessages(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := service.ListSessionMessages(uint(sessionID), page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// AskQuestion 智能问答接口，支持多轮对话
func AskQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.Ask(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

// ListAskHistory 获取用户的问答历史记录
func ListAskHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	conversationID, _ := strconv.Atoi(c.Query("conversation_id"))
	list, total, err := service.ListAskHistory(userID, page, size, uint(conversationID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}
