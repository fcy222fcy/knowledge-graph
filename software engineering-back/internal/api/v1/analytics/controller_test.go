package analytics

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockAnalyticsService 用于测试的 Mock 分析服务
type mockAnalyticsService struct {
	overviewResp   *response.OverviewResponse
	overviewErr    error
	hotPointsResp  []response.HotKnowledgePoint
	hotPointsErr   error
	masteryResp    []response.KnowledgeMastery
	masteryErr     error
	weakPointsResp []response.WeakPoint
	weakPointsErr  error
	trendsResp     *response.TrendData
	trendsErr      error
}

func (m *mockAnalyticsService) GetOverview(userID uint) (*response.OverviewResponse, error) {
	return m.overviewResp, m.overviewErr
}

func (m *mockAnalyticsService) GetHotKnowledgePoints(limit int) ([]response.HotKnowledgePoint, error) {
	return m.hotPointsResp, m.hotPointsErr
}

func (m *mockAnalyticsService) GetKnowledgeMastery(userID uint) ([]response.KnowledgeMastery, error) {
	return m.masteryResp, m.masteryErr
}

func (m *mockAnalyticsService) GetWeakPoints(userID uint, limit int) ([]response.WeakPoint, error) {
	return m.weakPointsResp, m.weakPointsErr
}

func (m *mockAnalyticsService) GetTrends(userID uint, days int) (*response.TrendData, error) {
	return m.trendsResp, m.trendsErr
}

func setupTestRouter(ctrl *AnalyticsController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.GET("/analytics/overview", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.GetOverview(c)
		})
		api.GET("/analytics/hot-knowledge-points", ctrl.GetHotKnowledgePoints)
		api.GET("/analytics/knowledge-mastery", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.GetKnowledgeMastery(c)
		})
		api.GET("/analytics/weak-points", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.GetWeakPoints(c)
		})
		api.GET("/analytics/trends", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.GetTrends(c)
		})
	}
	return r
}

func TestAnalyticsController_GetOverview_Success(t *testing.T) {
	mockSvc := &mockAnalyticsService{
		overviewResp: &response.OverviewResponse{
			TodayQuestionsAsked:     5,
			TotalQuestionsAsked:     100,
			TotalQuizzesTaken:       50,
			AverageCorrectRate:      0.75,
			KnowledgePointsMastered: 10,
			KnowledgePointsTotal:    20,
			MasteryRate:             0.5,
		},
	}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/overview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["today_questions_asked"] != float64(5) {
		t.Errorf("today_questions_asked = %v, want 5", data["today_questions_asked"])
	}
}

func TestAnalyticsController_GetOverview_Error(t *testing.T) {
	mockSvc := &mockAnalyticsService{overviewErr: errors.New("查询失败")}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/overview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestAnalyticsController_GetHotKnowledgePoints_Success(t *testing.T) {
	mockSvc := &mockAnalyticsService{
		hotPointsResp: []response.HotKnowledgePoint{
			{KnowledgePointID: 1, KnowledgePointName: "软件工程", Heat: 100},
			{KnowledgePointID: 2, KnowledgePointName: "需求分析", Heat: 80},
		},
	}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/hot-knowledge-points?limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("数据长度 = %v, want 2", len(data))
	}
}

func TestAnalyticsController_GetKnowledgeMastery_Success(t *testing.T) {
	mockSvc := &mockAnalyticsService{
		masteryResp: []response.KnowledgeMastery{
			{KnowledgePointID: 1, KnowledgePointName: "软件工程", MasteryRate: 0.9, Level: "mastered"},
			{KnowledgePointID: 2, KnowledgePointName: "需求分析", MasteryRate: 0.5, Level: "learning"},
		},
	}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/knowledge-mastery", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("数据长度 = %v, want 2", len(data))
	}
}

func TestAnalyticsController_GetWeakPoints_Success(t *testing.T) {
	mockSvc := &mockAnalyticsService{
		weakPointsResp: []response.WeakPoint{
			{KnowledgePointID: 1, KnowledgePointName: "设计模式", CorrectRate: 0.3},
		},
	}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/weak-points?limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("数据长度 = %v, want 1", len(data))
	}
}

func TestAnalyticsController_GetTrends_Success(t *testing.T) {
	mockSvc := &mockAnalyticsService{
		trendsResp: &response.TrendData{
			DailyStats: []response.DailyStat{
				{Date: "2024-01-01", QuestionsAsked: 10, CorrectRate: 0.8},
				{Date: "2024-01-02", QuestionsAsked: 15, CorrectRate: 0.7},
			},
		},
	}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/trends?days=7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	dailyStats := data["daily_stats"].([]interface{})
	if len(dailyStats) != 2 {
		t.Errorf("daily_stats 长度 = %v, want 2", len(dailyStats))
	}
}

func TestAnalyticsController_GetTrends_Error(t *testing.T) {
	mockSvc := &mockAnalyticsService{trendsErr: errors.New("查询失败")}
	ctrl := NewAnalyticsController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/analytics/trends?days=7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
