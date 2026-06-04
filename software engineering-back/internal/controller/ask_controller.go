package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func CreateSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.CreateSession(userID, req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := service.ListSessions(userID, page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}

func ListSessionMessages(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := service.ListSessionMessages(uint(sessionID), page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}

func AskQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.Ask(userID, req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListAskHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	conversationID, _ := strconv.Atoi(c.Query("conversation_id"))
	list, total, err := service.ListAskHistory(userID, page, size, uint(conversationID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
