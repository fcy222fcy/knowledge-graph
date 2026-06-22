package document

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

// mockDocumentService 用于测试的 Mock 文档服务
type mockDocumentService struct {
	uploadResp  *response.DocumentResponse
	uploadErr   error
	getResp     *response.DocumentResponse
	getErr      error
	contentResp *response.DocumentContentResponse
	contentErr  error
	updateErr   error
	deleteErr   error
	listResp    []response.DocumentResponse
	listTotal   int64
	listErr     error
}

func (m *mockDocumentService) UploadDocument(userID uint, title, description, filename string, fileSize int64, fileType string, reader io.Reader) (*response.DocumentResponse, error) {
	return m.uploadResp, m.uploadErr
}

func (m *mockDocumentService) GetDocument(id uint) (*response.DocumentResponse, error) {
	return m.getResp, m.getErr
}

func (m *mockDocumentService) GetDocumentContent(id uint) (*response.DocumentContentResponse, error) {
	return m.contentResp, m.contentErr
}

func (m *mockDocumentService) UpdateDocument(id uint, req request.UpdateDocumentRequest) error {
	return m.updateErr
}

func (m *mockDocumentService) DeleteDocument(userID uint, id uint) error {
	return m.deleteErr
}

func (m *mockDocumentService) ListUserDocuments(userID uint, page, size int, keyword, status string) ([]response.DocumentResponse, int64, error) {
	return m.listResp, m.listTotal, m.listErr
}

func setupTestRouter(ctrl *DocumentController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.POST("/documents/upload", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.UploadDocument(c)
		})
		api.GET("/documents/:id", ctrl.GetDocument)
		api.GET("/documents/:id/content", ctrl.GetDocumentContent)
		api.PUT("/documents/:id", ctrl.UpdateDocument)
		api.DELETE("/documents/:id", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.DeleteDocument(c)
		})
		api.GET("/documents", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.ListDocuments(c)
		})
	}
	return r
}

func TestDocumentController_GetDocument_Success(t *testing.T) {
	mockSvc := &mockDocumentService{
		getResp: &response.DocumentResponse{
			ID:       1,
			Title:    "测试文档",
			Filename: "test.md",
			FileSize: 1024,
			FileType: ".md",
			Status:   "completed",
		},
	}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/documents/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["title"] != "测试文档" {
		t.Errorf("title = %v, want '测试文档'", data["title"])
	}
}

func TestDocumentController_GetDocument_NotFound(t *testing.T) {
	mockSvc := &mockDocumentService{getErr: errors.New("文档不存在")}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/documents/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestDocumentController_GetDocumentContent_Success(t *testing.T) {
	mockSvc := &mockDocumentService{
		contentResp: &response.DocumentContentResponse{
			ID:      1,
			Title:   "测试文档",
			Content: "# 测试文档\n\n这是一个测试文档。",
		},
	}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/documents/1/content", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["content"] == nil {
		t.Error("content 不应为空")
	}
}

func TestDocumentController_UpdateDocument_Success(t *testing.T) {
	mockSvc := &mockDocumentService{}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.UpdateDocumentRequest{
		Title:       "更新后的标题",
		Description: "更新后的描述",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/documents/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestDocumentController_UpdateDocument_InvalidJSON(t *testing.T) {
	mockSvc := &mockDocumentService{}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("PUT", "/api/v1/documents/1", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestDocumentController_DeleteDocument_Success(t *testing.T) {
	mockSvc := &mockDocumentService{}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("DELETE", "/api/v1/documents/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestDocumentController_DeleteDocument_Error(t *testing.T) {
	mockSvc := &mockDocumentService{deleteErr: errors.New("删除失败")}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("DELETE", "/api/v1/documents/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestDocumentController_ListDocuments_Success(t *testing.T) {
	mockSvc := &mockDocumentService{
		listResp: []response.DocumentResponse{
			{ID: 1, Title: "文档1"},
			{ID: 2, Title: "文档2"},
		},
		listTotal: 2,
	}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/documents?page=1&size=10", nil)
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

func TestDocumentController_ListDocuments_Error(t *testing.T) {
	mockSvc := &mockDocumentService{listErr: errors.New("查询失败")}
	ctrl := NewDocumentController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/documents?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestFixFilenameEncodingIface_ASCII(t *testing.T) {
	result := fixFilenameEncodingIface("test.md")
	if result != "test.md" {
		t.Errorf("fixFilenameEncodingIface() = %v, want 'test.md'", result)
	}
}

func TestFixFilenameEncodingIface_Empty(t *testing.T) {
	result := fixFilenameEncodingIface("")
	if result != "" {
		t.Errorf("fixFilenameEncodingIface() = %v, want ''", result)
	}
}
