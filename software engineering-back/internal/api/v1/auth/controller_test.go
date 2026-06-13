package auth

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

// mockAuthService 用于测试的 Mock 认证服务
type mockAuthService struct {
	registerErr   error
	loginResp     *response.LoginResponse
	loginErr      error
	refreshToken  string
	refreshErr    error
}

func (m *mockAuthService) Register(req request.RegisterRequest) error {
	return m.registerErr
}

func (m *mockAuthService) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	return m.loginResp, m.loginErr
}

func (m *mockAuthService) RefreshToken(oldToken string) (string, error) {
	return m.refreshToken, m.refreshErr
}

func setupTestRouter(ctrl *AuthController) *gin.Engine {
	r := gin.New()
	api := r.Group("/api/v1/auth")
	{
		api.POST("/register", ctrl.Register)
		api.POST("/login", ctrl.Login)
		api.POST("/refresh", ctrl.Refresh)
	}
	return r
}

func TestAuthController_Register_Success(t *testing.T) {
	mockSvc := &mockAuthService{}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
		Nickname: "测试用户",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
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
}

func TestAuthController_Register_Error(t *testing.T) {
	mockSvc := &mockAuthService{
		registerErr: errors.New("用户名已存在"),
	}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.RegisterRequest{
		Username: "existinguser",
		Password: "password123",
		Email:    "test@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "用户名已存在" {
		t.Errorf("message = %v, want '用户名已存在'", resp["message"])
	}
}

func TestAuthController_Login_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		loginResp: &response.LoginResponse{
			Token: "test-token-123",
			User: response.UserResponse{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
	}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data := resp["data"].(map[string]interface{})
	if data["token"] != "test-token-123" {
		t.Errorf("token = %v, want 'test-token-123'", data["token"])
	}
}

func TestAuthController_Login_Unauthorized(t *testing.T) {
	mockSvc := &mockAuthService{
		loginErr: errors.New("用户名或密码错误"),
	}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.LoginRequest{
		Username: "wronguser",
		Password: "wrongpass",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "用户名或密码错误" {
		t.Errorf("message = %v, want '用户名或密码错误'", resp["message"])
	}
}

func TestAuthController_Refresh_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		refreshToken: "new-token-456",
	}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.RefreshRequest{
		Token: "old-token-123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data := resp["data"].(map[string]interface{})
	if data["token"] != "new-token-456" {
		t.Errorf("token = %v, want 'new-token-456'", data["token"])
	}
}

func TestAuthController_Refresh_Unauthorized(t *testing.T) {
	mockSvc := &mockAuthService{
		refreshErr: errors.New("无效的令牌"),
	}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	body := request.RefreshRequest{
		Token: "invalid-token",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthController_Register_InvalidJSON(t *testing.T) {
	mockSvc := &mockAuthService{}
	ctrl := NewAuthController(mockSvc)
	r := setupTestRouter(ctrl)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusBadRequest)
	}
}
