package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"software_engineering/pkg/jwt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(RequireAuth())

	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"user_id": userID,
			"username": username,
		})
	})

	return r
}

func TestRequireAuth_NoHeader(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "未授权" {
		t.Errorf("message = %v, want '未授权'", resp["message"])
	}
}

func TestRequireAuth_InvalidFormat(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		name   string
		header string
	}{
		{"无Bearer前缀", "token123"},
		{"只有Bearer", "Bearer"},
		{"Bearer后无空格", "Bearertoken123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("状态码 = %v, want %v", w.Code, http.StatusUnauthorized)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["message"] != "无效的令牌格式" {
				t.Errorf("message = %v, want '无效的令牌格式'", resp["message"])
			}
		})
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusUnauthorized)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "无效或过期的令牌" {
		t.Errorf("message = %v, want '无效或过期的令牌'", resp["message"])
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	r := setupRouter()

	// 生成有效 token
	token, err := jwt.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["user_id"] != float64(1) {
		t.Errorf("user_id = %v, want 1", resp["user_id"])
	}
	if resp["username"] != "testuser" {
		t.Errorf("username = %v, want 'testuser'", resp["username"])
	}
}

func TestRequireAuth_EmptyToken(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}
