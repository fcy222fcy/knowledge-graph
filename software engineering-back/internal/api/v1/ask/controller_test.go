package ask

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockAskService 用于测试的 Mock 问答服务
type mockAskService struct {
	createSessionResp *response.AskSessionResponse
	createSessionErr  error
	sessionsResp      []response.AskSessionResponse
	sessionsTotal     int64
	sessionsErr       error
	messagesResp      []response.AskMessageResponse
	messagesTotal     int64
	messagesErr       error
	askResp           *response.AskResponse
	askErr            error
	historyResp       []response.AskHistoryItem
	historyTotal      int64
	historyErr        error
}

func (m *mockAskService) CreateSession(userID uint, req request.CreateSessionRequest) (*response.AskSessionResponse, error) {
	return m.createSessionResp, m.createSessionErr
}

func (m *mockAskService) ListSessions(userID uint, page, size int) ([]response.AskSessionResponse, int64, error) {
	return m.sessionsResp, m.sessionsTotal, m.sessionsErr
}

func (m *mockAskService) ListSessionMessages(sessionID uint, page, size int) ([]response.AskMessageResponse, int64, error) {
	return m.messagesResp, m.messagesTotal, m.messagesErr
}

func (m *mockAskService) Ask(userID uint, req request.AskRequest) (*response.AskResponse, error) {
	return m.askResp, m.askErr
}

func (m *mockAskService) ListAskHistory(userID uint, page, size int, conversationID uint) ([]response.AskHistoryItem, int64, error) {
	return m.historyResp, m.historyTotal, m.historyErr
}

func setupTestRouter(ctrl *AskController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.POST("/ask/sessions", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.CreateSession(c)
		})
		api.GET("/ask/sessions", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.ListSessions(c)
		})
		api.GET("/ask/sessions/:id/messages", ctrl.ListSessionMessages)
		api.POST("/ask/question", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.AskQuestion(c)
		})
		api.GET("/ask/history", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.ListAskHistory(c)
		})
	}
	return r
}

func TestAskController_CreateSession_Success(t *testing.T) {
	mockSvc := &mockAskService{
		createSessionResp: &response.AskSessionResponse{
			ConversationID: 1,
			Title:          "新会话",
			MessageCount:   0,
			UpdatedAt:      "2024-01-01T00:00:00Z",
		},
	}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.CreateSessionRequest{Title: "新会话"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/ask/sessions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["title"] != "新会话" {
		t.Errorf("title = %v, want '新会话'", data["title"])
	}
}

func TestAskController_CreateSession_InvalidJSON(t *testing.T) {
	mockSvc := &mockAskService{}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/ask/sessions", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestAskController_ListSessions_Success(t *testing.T) {
	mockSvc := &mockAskService{
		sessionsResp: []response.AskSessionResponse{
			{ConversationID: 1, Title: "会话1"},
			{ConversationID: 2, Title: "会话2"},
		},
		sessionsTotal: 2,
	}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/ask/sessions?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(2) {
		t.Errorf("total = %v, want 2", data["total"])
	}
}

func TestAskController_ListSessionMessages_Success(t *testing.T) {
	mockSvc := &mockAskService{
		messagesResp: []response.AskMessageResponse{
			{MessageID: 1, Role: "user", Content: "什么是软件工程？"},
			{MessageID: 2, Role: "assistant", Content: "软件工程是一门学科"},
		},
		messagesTotal: 2,
	}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/ask/sessions/1/messages?page=1&size=20", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(2) {
		t.Errorf("total = %v, want 2", data["total"])
	}
}

func TestAskController_AskQuestion_Success(t *testing.T) {
	mockSvc := &mockAskService{
		askResp: &response.AskResponse{
			ConversationID: 1,
			QuestionID:     1,
			Answer:         "软件工程是一门系统化、规范化、可量化的方法应用于软件开发的学科。",
			Confidence:     0.85,
			Sources: []response.AskSource{
				{DocumentID: 1, DocumentTitle: "软件工程导论", Content: "软件工程是..."},
			},
			CreatedAt: "2024-01-01T00:00:00Z",
		},
	}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.AskRequest{Question: "什么是软件工程？"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/ask/question", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["answer"] == nil || data["answer"] == "" {
		t.Error("answer 不应为空")
	}
}

func TestAskController_AskQuestion_Error(t *testing.T) {
	mockSvc := &mockAskService{askErr: errors.New("服务不可用")}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.AskRequest{Question: "测试问题"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/ask/question", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestAskController_AskQuestion_InvalidJSON(t *testing.T) {
	mockSvc := &mockAskService{}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/ask/question", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestAskController_ListAskHistory_Success(t *testing.T) {
	mockSvc := &mockAskService{
		historyResp: []response.AskHistoryItem{
			{ConversationID: 1, Title: "会话1", LastQuestion: "什么是软件工程？"},
			{ConversationID: 2, Title: "会话2", LastQuestion: "什么是需求分析？"},
		},
		historyTotal: 2,
	}
	ctrl := NewAskController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/ask/history?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(2) {
		t.Errorf("total = %v, want 2", data["total"])
	}
}
