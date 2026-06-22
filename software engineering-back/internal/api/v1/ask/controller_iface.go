package ask

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// AskService 定义问答服务接口
type AskService interface {
	CreateSession(userID uint, req request.CreateSessionRequest) (*response.AskSessionResponse, error)            // 创建会话
	ListSessions(userID uint, page, size int) ([]response.AskSessionResponse, int64, error)                     // 列表会话
	ListSessionMessages(sessionID uint, page, size int) ([]response.AskMessageResponse, int64, error)            // 列表会话消息
	Ask(userID uint, req request.AskRequest) (*response.AskResponse, error)                                      // 智能问答
	ListAskHistory(userID uint, page, size int, conversationID uint) ([]response.AskHistoryItem, int64, error)   // 列表问答历史
}

// AskController 问答控制器
type AskController struct {
	askService AskService // 问答服务
}

// NewAskController 创建问答控制器实例
func NewAskController(askService AskService) *AskController {
	return &AskController{askService: askService}
}

// CreateSession 创建问答会话
func (ctrl *AskController) CreateSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := ctrl.askService.CreateSession(userID, req)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// ListSessions 获取问答会话列表
func (ctrl *AskController) ListSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := ctrl.askService.ListSessions(userID, page, size)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}

// ListSessionMessages 获取会话消息记录
func (ctrl *AskController) ListSessionMessages(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := ctrl.askService.ListSessionMessages(uint(sessionID), page, size)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}

// AskQuestion 智能问答
func (ctrl *AskController) AskQuestion(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := ctrl.askService.Ask(userID, req)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// ListAskHistory 获取问答历史记录
func (ctrl *AskController) ListAskHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	conversationID, _ := strconv.Atoi(c.Query("conversation_id"))
	list, total, err := ctrl.askService.ListAskHistory(userID, page, size, uint(conversationID))
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}
