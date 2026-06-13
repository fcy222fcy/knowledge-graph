package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthCheck(t *testing.T) {
	r := gin.New()
	r.GET("/api/v1/health", healthCheck)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp["code"] != float64(200) {
		t.Errorf("code = %v, want 200", resp["code"])
	}

	data := resp["data"].(map[string]interface{})
	if data["status"] != "ok" {
		t.Errorf("status = %v, want 'ok'", data["status"])
	}
	if data["service"] != "software-engineering-backend" {
		t.Errorf("service = %v, want 'software-engineering-backend'", data["service"])
	}
}

func TestSetupRoutes_HealthEndpoint(t *testing.T) {
	r := gin.New()
	SetupRoutes(r)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("健康检查端点状态码 = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestSetupRoutes_ProtectedEndpointsRequireAuth(t *testing.T) {
	r := gin.New()
	SetupRoutes(r)

	// 测试需要认证的端点（使用已知存在的路由）
	protectedEndpoints := []string{
		"/api/v1/ask/history",
		"/api/v1/analytics/overview",
	}

	for _, endpoint := range protectedEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", endpoint, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// 不带 token 应该返回 401
			if w.Code != http.StatusUnauthorized {
				t.Errorf("GET %v 状态码 = %v, want %v", endpoint, w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestSetupRoutes_PublicEndpoints(t *testing.T) {
	r := gin.New()
	SetupRoutes(r)

	// 测试公开端点
	publicEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/health"},
		{"POST", "/api/v1/auth/register"},
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/auth/refresh"},
	}

	for _, ep := range publicEndpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			req := httptest.NewRequest(ep.method, ep.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// 公开端点不应该返回 401（可能返回其他错误如参数缺失，但不是 401）
			if w.Code == http.StatusUnauthorized {
				t.Errorf("公开端点 %v %v 不应该返回 401", ep.method, ep.path)
			}
		})
	}
}
