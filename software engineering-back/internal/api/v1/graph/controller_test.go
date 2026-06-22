package graph

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

// mockGraphService 用于测试的 Mock 图谱服务
type mockGraphService struct {
	graphData   *response.GraphDataResponse
	graphErr    error
	buildResp   *response.BuildGraphResponse
	buildErr    error
	latestResp  *response.BuildGraphResponse
	latestErr   error
	historyResp *response.BuildHistoryResponse
	historyErr  error
}

func (m *mockGraphService) GetGraphData(documentID uint, keyword, relationType string) (*response.GraphDataResponse, error) {
	return m.graphData, m.graphErr
}

func (m *mockGraphService) BuildGraph(documentIDs []uint) (*response.BuildGraphResponse, error) {
	return m.buildResp, m.buildErr
}

func (m *mockGraphService) GetLatestBuildResult() (*response.BuildGraphResponse, error) {
	return m.latestResp, m.latestErr
}

func (m *mockGraphService) ListBuildHistory(page, size int) (*response.BuildHistoryResponse, error) {
	return m.historyResp, m.historyErr
}

func setupTestRouter(ctrl *GraphController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.GET("/graph", ctrl.GetGraph)
		api.POST("/graph/build", ctrl.BuildGraph)
		api.GET("/graph/build/latest", ctrl.GetLatestBuild)
		api.GET("/graph/build/history", ctrl.ListBuildHistory)
	}
	return r
}

func TestGraphController_GetGraph_Success(t *testing.T) {
	mockSvc := &mockGraphService{
		graphData: &response.GraphDataResponse{
			Nodes: []response.GraphNode{
				{ID: 1, Name: "知识点1", Description: "描述1"},
				{ID: 2, Name: "知识点2", Description: "描述2"},
			},
			Edges: []response.GraphEdge{
				{ID: 1, Source: 1, Target: 2, RelationType: "RELATED"},
			},
			Summary: response.GraphSummary{NodeCount: 2, EdgeCount: 1},
		},
	}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/graph?document_id=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	nodes := data["nodes"].([]interface{})
	if len(nodes) != 2 {
		t.Errorf("nodes 长度 = %v, want 2", len(nodes))
	}
}

func TestGraphController_GetGraph_Error(t *testing.T) {
	mockSvc := &mockGraphService{graphErr: errors.New("查询失败")}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/graph", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGraphController_BuildGraph_Success(t *testing.T) {
	mockSvc := &mockGraphService{
		buildResp: &response.BuildGraphResponse{
			BuildID:          1,
			CreatedPoints:    5,
			CreatedRelations: 3,
			ChunkCount:       10,
			VectorCount:      10,
			Status:           "completed",
			Message:          "构建成功",
		},
	}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.BuildGraphRequest{DocumentIDs: []uint{1, 2}}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/graph/build", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["created_points"] != float64(5) {
		t.Errorf("created_points = %v, want 5", data["created_points"])
	}
}

func TestGraphController_BuildGraph_InvalidJSON(t *testing.T) {
	mockSvc := &mockGraphService{}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/graph/build", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestGraphController_BuildGraph_Error(t *testing.T) {
	mockSvc := &mockGraphService{buildErr: errors.New("构建失败")}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.BuildGraphRequest{DocumentIDs: []uint{1}}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/graph/build", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGraphController_GetLatestBuild_Success(t *testing.T) {
	mockSvc := &mockGraphService{
		latestResp: &response.BuildGraphResponse{
			BuildID:          1,
			CreatedPoints:    5,
			CreatedRelations: 3,
			Status:           "completed",
		},
	}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/graph/build/latest", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestGraphController_GetLatestBuild_NotFound(t *testing.T) {
	mockSvc := &mockGraphService{latestErr: errors.New("暂无构建记录")}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/graph/build/latest", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestGraphController_ListBuildHistory_Success(t *testing.T) {
	mockSvc := &mockGraphService{
		historyResp: &response.BuildHistoryResponse{
			List: []response.BuildGraphResponse{
				{BuildID: 1, Status: "completed"},
				{BuildID: 2, Status: "completed"},
			},
			Total:     2,
			Page:      1,
			Size:      10,
			TotalPage: 1,
		},
	}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/graph/build/history?page=1&size=10", nil)
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

func TestGraphController_ListBuildHistory_Error(t *testing.T) {
	mockSvc := &mockGraphService{historyErr: errors.New("查询失败")}
	ctrl := NewGraphController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/graph/build/history?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
