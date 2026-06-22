package user

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

// mockUserService 用于测试的 Mock 用户服务
type mockUserService struct {
	profileResp  *response.UserResponse
	profileErr   error
	updateErr    error
	changePwdErr error
	listResp     *response.UserListResponse
	listErr      error
}

func (m *mockUserService) GetProfile(userID uint) (*response.UserResponse, error) {
	return m.profileResp, m.profileErr
}

func (m *mockUserService) UpdateProfile(userID uint, req request.UpdateProfileRequest) error {
	return m.updateErr
}

func (m *mockUserService) ChangePassword(userID uint, req request.ChangePasswordRequest) error {
	return m.changePwdErr
}

func (m *mockUserService) ListUsers(page, size int) (*response.UserListResponse, error) {
	return m.listResp, m.listErr
}

func setupTestRouter(ctrl *UserController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1")
	{
		api.GET("/users/profile", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.GetProfile(c)
		})
		api.PUT("/users/profile", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.UpdateProfile(c)
		})
		api.POST("/users/password", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			ctrl.ChangePassword(c)
		})
		api.GET("/users/list", ctrl.ListUsers)
	}
	return r
}

func TestUserController_GetProfile_Success(t *testing.T) {
	mockSvc := &mockUserService{
		profileResp: &response.UserResponse{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Nickname: "测试用户",
			Status:   1,
		},
	}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["username"] != "testuser" {
		t.Errorf("username = %v, want 'testuser'", data["username"])
	}
}

func TestUserController_GetProfile_NotFound(t *testing.T) {
	mockSvc := &mockUserService{profileErr: errors.New("用户不存在")}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestUserController_UpdateProfile_Success(t *testing.T) {
	mockSvc := &mockUserService{}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.UpdateProfileRequest{
		Nickname: "新昵称",
		Avatar:   "https://example.com/avatar.jpg",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestUserController_UpdateProfile_InvalidJSON(t *testing.T) {
	mockSvc := &mockUserService{}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestUserController_UpdateProfile_Error(t *testing.T) {
	mockSvc := &mockUserService{updateErr: errors.New("更新失败")}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.UpdateProfileRequest{Nickname: "新昵称"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestUserController_ChangePassword_Success(t *testing.T) {
	mockSvc := &mockUserService{}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.ChangePasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/users/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestUserController_ChangePassword_InvalidJSON(t *testing.T) {
	mockSvc := &mockUserService{}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/users/password", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestUserController_ChangePassword_Error(t *testing.T) {
	mockSvc := &mockUserService{changePwdErr: errors.New("原密码错误")}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/users/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestUserController_ListUsers_Success(t *testing.T) {
	mockSvc := &mockUserService{
		listResp: &response.UserListResponse{
			List: []response.UserResponse{
				{ID: 1, Username: "user1"},
				{ID: 2, Username: "user2"},
			},
			Total:     2,
			Page:      1,
			Size:      10,
			TotalPage: 1,
		},
	}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/users/list?page=1&size=10", nil)
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

func TestUserController_ListUsers_Error(t *testing.T) {
	mockSvc := &mockUserService{listErr: errors.New("查询失败")}
	ctrl := NewUserController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("GET", "/api/v1/users/list?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
