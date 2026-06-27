package ask

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dto "software_engineering/internal/model/dto/response"
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

// AskQuestionStream 流式智能问答接口，使用 Server-Sent Events (SSE)
func AskQuestionStream(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 调用流式服务（在设置 SSE 头之前）
	sessionID, ch, err := service.AskStream(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Writer.Flush()

	// 先发送会话 ID
	c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", mustJSON(map[string]interface{}{
		"type":        "session",
		"session_id":  sessionID,
	}))))
	c.Writer.Flush()

	// 逐个发送事件
	for event := range ch {
		data := map[string]interface{}{
			"type":    event.Type,
			"content": event.Content,
		}
		if event.Type == "done" {
			data["confidence"] = event.Confidence
			if event.Sources == nil {
				event.Sources = make([]dto.AskSource, 0)
			}
			if event.Related == nil {
				event.Related = make([]dto.KPRef, 0)
			}
			data["sources"] = event.Sources
			data["related"] = event.Related
		}
		jsonData, _ := json.Marshal(data)
		c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonData)))
		c.Writer.Flush()
	}

	// 发送结束标记
	c.Writer.Write([]byte("data: [DONE]\n\n"))
	c.Writer.Flush()
}

// mustJSON JSON 序列化辅助函数
func mustJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
