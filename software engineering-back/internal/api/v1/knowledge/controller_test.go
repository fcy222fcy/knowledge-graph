package knowledge

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

// mockKnowledgeService 用于测试的 Mock 知识点服务
type mockKnowledgeService struct {
	createPointID   uint
	createPointErr  error
	getPointResp    *response.KnowledgePointResponse
	getPointErr     error
	updatePointErr  error
	deletePointErr  error
	listPointsResp  []response.KnowledgePointResponse
	listPointsTotal int64
	listPointsErr   error
	createRelID     uint
	createRelErr    error
	updateRelErr    error
	deleteRelErr    error
	listRelsResp    []response.KnowledgeRelationResponse
	listRelsTotal   int64
	listRelsErr     error
}

func (m *mockKnowledgeService) CreateKnowledgePoint(req request.CreateKnowledgePointRequest) (uint, error) {
	return m.createPointID, m.createPointErr
}

func (m *mockKnowledgeService) GetKnowledgePoint(id uint) (*response.KnowledgePointResponse, error) {
	return m.getPointResp, m.getPointErr
}

func (m *mockKnowledgeService) UpdateKnowledgePoint(id uint, req request.UpdateKnowledgePointRequest) error {
	return m.updatePointErr
}

func (m *mockKnowledgeService) DeleteKnowledgePoint(id uint) error {
	return m.deletePointErr
}

func (m *mockKnowledgeService) ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]response.KnowledgePointResponse, int64, error) {
	return m.listPointsResp, m.listPointsTotal, m.listPointsErr
}

func (m *mockKnowledgeService) CreateRelation(req request.CreateRelationRequest) (uint, error) {
	return m.createRelID, m.createRelErr
}

func (m *mockKnowledgeService) UpdateRelation(id uint, req request.UpdateRelationRequest) error {
	return m.updateRelErr
}

func (m *mockKnowledgeService) DeleteRelation(id uint) error {
	return m.deleteRelErr
}

func (m *mockKnowledgeService) ListRelations(page, size int, pointID uint) ([]response.KnowledgeRelationResponse, int64, error) {
	return m.listRelsResp, m.listRelsTotal, m.listRelsErr
}

func setupTestRouter(ctrl *KnowledgeController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.POST("/knowledge/points", ctrl.CreateKnowledgePoint)
		api.GET("/knowledge/points/:id", ctrl.GetKnowledgePoint)
		api.PUT("/knowledge/points/:id", ctrl.UpdateKnowledgePoint)
		api.DELETE("/knowledge/points/:id", ctrl.DeleteKnowledgePoint)
		api.GET("/knowledge/points", ctrl.ListKnowledgePoints)
		api.POST("/knowledge/relations", ctrl.CreateRelation)
		api.PUT("/knowledge/relations/:id", ctrl.UpdateRelation)
		api.DELETE("/knowledge/relations/:id", ctrl.DeleteRelation)
		api.GET("/knowledge/relations", ctrl.ListRelations)
	}
	return r
}

func TestKnowledgeController_CreateKnowledgePoint_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{createPointID: 1}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.CreateKnowledgePointRequest{
		Name:        "软件工程",
		Description: "软件工程是一门学科",
		DocumentID:  1,
		Category:    "概念",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/knowledge/points", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["id"] != float64(1) {
		t.Errorf("id = %v, want 1", data["id"])
	}
}

func TestKnowledgeController_CreateKnowledgePoint_Error(t *testing.T) {
	mockSvc := &mockKnowledgeService{createPointErr: errors.New("创建失败")}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.CreateKnowledgePointRequest{Name: "测试"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/knowledge/points", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestKnowledgeController_GetKnowledgePoint_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{
		getPointResp: &response.KnowledgePointResponse{
			ID:          1,
			Name:        "软件工程",
			Description: "软件工程是一门学科",
			DocumentID:  1,
			Category:    "概念",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
	}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/knowledge/points/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["name"] != "软件工程" {
		t.Errorf("name = %v, want '软件工程'", data["name"])
	}
}

func TestKnowledgeController_GetKnowledgePoint_NotFound(t *testing.T) {
	mockSvc := &mockKnowledgeService{getPointErr: errors.New("知识点不存在")}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/knowledge/points/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestKnowledgeController_UpdateKnowledgePoint_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.UpdateKnowledgePointRequest{Name: "更新后的知识点"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/knowledge/points/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestKnowledgeController_DeleteKnowledgePoint_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("DELETE", "/api/v1/knowledge/points/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestKnowledgeController_ListKnowledgePoints_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{
		listPointsResp: []response.KnowledgePointResponse{
			{ID: 1, Name: "知识点1"},
			{ID: 2, Name: "知识点2"},
		},
		listPointsTotal: 2,
	}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/knowledge/points?page=1&size=10", nil)
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

func TestKnowledgeController_CreateRelation_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{createRelID: 1}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.CreateRelationRequest{
		SourceID:     1,
		TargetID:     2,
		RelationType: "RELATED",
		Description:  "相关关系",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/knowledge/relations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestKnowledgeController_CreateRelation_InvalidJSON(t *testing.T) {
	mockSvc := &mockKnowledgeService{}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/knowledge/relations", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestKnowledgeController_DeleteRelation_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("DELETE", "/api/v1/knowledge/relations/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestKnowledgeController_ListRelations_Success(t *testing.T) {
	mockSvc := &mockKnowledgeService{
		listRelsResp: []response.KnowledgeRelationResponse{
			{ID: 1, SourceID: 1, TargetID: 2, RelationType: "RELATED"},
		},
		listRelsTotal: 1,
	}
	ctrl := NewKnowledgeController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/knowledge/relations?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(1) {
		t.Errorf("total = %v, want 1", data["total"])
	}
}
