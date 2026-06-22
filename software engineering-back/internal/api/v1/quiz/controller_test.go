package quiz

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

// mockQuizService 用于测试的 Mock 答题服务
type mockQuizService struct {
	submitResp *response.QuizResponse
	submitErr  error
	detailResp *response.QuizResponse
	detailErr  error
	historyResp []response.QuizResponse
	historyTotal int64
	historyErr  error
}

func (m *mockQuizService) SubmitQuiz(userID uint, req request.SubmitQuizRequest) (*response.QuizResponse, error) {
	return m.submitResp, m.submitErr
}

func (m *mockQuizService) GetQuizDetail(id uint) (*response.QuizResponse, error) {
	return m.detailResp, m.detailErr
}

func (m *mockQuizService) ListQuizHistory(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]response.QuizResponse, int64, error) {
	return m.historyResp, m.historyTotal, m.historyErr
}

func setupTestRouter(ctrl *QuizController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.POST("/quizzes/submit", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.SubmitQuiz(c)
		})
		api.GET("/quizzes/:id", ctrl.GetQuizDetail)
		api.GET("/quizzes/history", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.ListQuizHistory(c)
		})
	}
	return r
}

func TestQuizController_SubmitQuiz_Success(t *testing.T) {
	mockSvc := &mockQuizService{
		submitResp: &response.QuizResponse{
			QuizID:        1,
			QuestionID:    1,
			UserAnswer:    "A",
			CorrectAnswer: "A",
			IsCorrect:     true,
			Explanation:   "正确答案是A",
			CreatedAt:     "2024-01-01T00:00:00Z",
		},
	}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.SubmitQuizRequest{
		QuestionID: 1,
		UserAnswer: "A",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/quizzes/submit", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["is_correct"] != true {
		t.Errorf("is_correct = %v, want true", data["is_correct"])
	}
}

func TestQuizController_SubmitQuiz_Error(t *testing.T) {
	mockSvc := &mockQuizService{submitErr: errors.New("题目不存在")}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.SubmitQuizRequest{
		QuestionID: 999,
		UserAnswer: "A",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/quizzes/submit", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestQuizController_SubmitQuiz_InvalidJSON(t *testing.T) {
	mockSvc := &mockQuizService{}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/quizzes/submit", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestQuizController_GetQuizDetail_Success(t *testing.T) {
	mockSvc := &mockQuizService{
		detailResp: &response.QuizResponse{
			QuizID:        1,
			QuestionID:    1,
			QuestionTitle: "什么是软件工程？",
			Type:          "single",
			Difficulty:    "easy",
			Options: []response.QuestionOption{
				{Key: "A", Value: "软件开发方法论"},
				{Key: "B", Value: "一种编程语言"},
			},
			UserAnswer:    "A",
			CorrectAnswer: "A",
			IsCorrect:     true,
			Explanation:   "正确答案是A",
			CreatedAt:     "2024-01-01T00:00:00Z",
		},
	}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/quizzes/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["question_title"] != "什么是软件工程？" {
		t.Errorf("question_title = %v, want '什么是软件工程？'", data["question_title"])
	}
}

func TestQuizController_GetQuizDetail_NotFound(t *testing.T) {
	mockSvc := &mockQuizService{detailErr: errors.New("答题记录不存在")}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/quizzes/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestQuizController_ListQuizHistory_Success(t *testing.T) {
	mockSvc := &mockQuizService{
		historyResp: []response.QuizResponse{
			{QuizID: 1, QuestionID: 1, IsCorrect: true},
			{QuizID: 2, QuestionID: 2, IsCorrect: false},
		},
		historyTotal: 2,
	}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/quizzes/history?page=1&size=10", nil)
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

func TestQuizController_ListQuizHistory_Error(t *testing.T) {
	mockSvc := &mockQuizService{historyErr: errors.New("查询失败")}
	ctrl := NewQuizController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/quizzes/history?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
