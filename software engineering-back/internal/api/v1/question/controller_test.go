package question

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

// mockQuestionService 用于测试的 Mock 题目服务
type mockQuestionService struct {
	createErr     error
	createID      uint
	getResp       *response.QuestionResponse
	getErr        error
	updateErr     error
	deleteErr     error
	listResp      []response.QuestionResponse
	listTotal     int64
	listErr       error
}

func (m *mockQuestionService) CreateQuestion(req request.CreateQuestionRequest) (uint, error) {
	return m.createID, m.createErr
}

func (m *mockQuestionService) GetQuestion(id uint, includeAnswer bool) (*response.QuestionResponse, error) {
	return m.getResp, m.getErr
}

func (m *mockQuestionService) UpdateQuestion(id uint, req request.UpdateQuestionRequest) error {
	return m.updateErr
}

func (m *mockQuestionService) DeleteQuestion(id uint) error {
	return m.deleteErr
}

func (m *mockQuestionService) ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]response.QuestionResponse, int64, error) {
	return m.listResp, m.listTotal, m.listErr
}

func setupTestRouter(ctrl *QuestionController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.POST("/questions", ctrl.CreateQuestion)
		api.GET("/questions/:id", ctrl.GetQuestion)
		api.PUT("/questions/:id", ctrl.UpdateQuestion)
		api.DELETE("/questions/:id", ctrl.DeleteQuestion)
		api.GET("/questions", ctrl.ListQuestions)
	}
	return r
}

func TestQuestionController_CreateQuestion_Success(t *testing.T) {
	mockSvc := &mockQuestionService{createID: 1}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.CreateQuestionRequest{
		Title:            "什么是软件工程？",
		Type:             "single",
		Difficulty:       "easy",
		KnowledgePointID: 1,
		Options: []request.QuestionOption{
			{Key: "A", Value: "软件开发方法论"},
			{Key: "B", Value: "一种编程语言"},
		},
		Answer:      "A",
		Explanation: "软件工程是系统化、规范化、可量化的方法应用于软件开发。",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(200) {
		t.Errorf("code = %v, want 200", resp["code"])
	}

	data := resp["data"].(map[string]interface{})
	if data["id"] != float64(1) {
		t.Errorf("id = %v, want 1", data["id"])
	}
}

func TestQuestionController_CreateQuestion_Error(t *testing.T) {
	mockSvc := &mockQuestionService{createErr: errors.New("创建失败")}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.CreateQuestionRequest{
		Title:            "测试题目",
		Type:             "single",
		Difficulty:       "easy",
		KnowledgePointID: 1,
		Options: []request.QuestionOption{
			{Key: "A", Value: "选项A"},
			{Key: "B", Value: "选项B"},
		},
		Answer: "A",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/questions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestQuestionController_CreateQuestion_InvalidJSON(t *testing.T) {
	mockSvc := &mockQuestionService{}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/questions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestQuestionController_GetQuestion_Success(t *testing.T) {
	mockSvc := &mockQuestionService{
		getResp: &response.QuestionResponse{
			ID:               1,
			Title:            "什么是软件工程？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 1,
			Options: []response.QuestionOption{
				{Key: "A", Value: "软件开发方法论"},
				{Key: "B", Value: "一种编程语言"},
			},
			Answer:      "A",
			Explanation: "软件工程是系统化、规范化、可量化的方法应用于软件开发。",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
	}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/questions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["title"] != "什么是软件工程？" {
		t.Errorf("title = %v, want '什么是软件工程？'", data["title"])
	}
}

func TestQuestionController_GetQuestion_NotFound(t *testing.T) {
	mockSvc := &mockQuestionService{getErr: errors.New("题目不存在")}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/questions/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestQuestionController_UpdateQuestion_Success(t *testing.T) {
	mockSvc := &mockQuestionService{}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.UpdateQuestionRequest{
		Title: "更新后的题目",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/questions/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestQuestionController_UpdateQuestion_Error(t *testing.T) {
	mockSvc := &mockQuestionService{updateErr: errors.New("题目不存在")}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.UpdateQuestionRequest{Title: "更新"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/questions/999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestQuestionController_DeleteQuestion_Success(t *testing.T) {
	mockSvc := &mockQuestionService{}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("DELETE", "/api/v1/questions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestQuestionController_DeleteQuestion_Error(t *testing.T) {
	mockSvc := &mockQuestionService{deleteErr: errors.New("题目不存在")}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("DELETE", "/api/v1/questions/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestQuestionController_ListQuestions_Success(t *testing.T) {
	mockSvc := &mockQuestionService{
		listResp: []response.QuestionResponse{
			{ID: 1, Title: "题目1"},
			{ID: 2, Title: "题目2"},
		},
		listTotal: 2,
	}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/questions?page=1&size=10", nil)
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

func TestQuestionController_ListQuestions_Error(t *testing.T) {
	mockSvc := &mockQuestionService{listErr: errors.New("查询失败")}
	ctrl := NewQuestionController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/questions?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
