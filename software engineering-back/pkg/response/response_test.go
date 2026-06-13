package response

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

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccess(t *testing.T) {
	c, w := setupTestContext()

	Success(c, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("Success() status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp["code"] != float64(200) {
		t.Errorf("code = %v, want 200", resp["code"])
	}
	if resp["message"] != "success" {
		t.Errorf("message = %v, want 'success'", resp["message"])
	}
	if resp["data"] == nil {
		t.Error("data 不应为 nil")
	}
}

func TestSuccessWithMessage(t *testing.T) {
	c, w := setupTestContext()

	SuccessWithMessage(c, "自定义消息", "test-data")

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp["message"] != "自定义消息" {
		t.Errorf("message = %v, want '自定义消息'", resp["message"])
	}
	if resp["data"] != "test-data" {
		t.Errorf("data = %v, want 'test-data'", resp["data"])
	}
}

func TestError(t *testing.T) {
	c, w := setupTestContext()

	Error(c, http.StatusBadRequest, "请求参数错误")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Error() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if resp["code"] != float64(http.StatusBadRequest) {
		t.Errorf("code = %v, want %v", resp["code"], http.StatusBadRequest)
	}
	if resp["message"] != "请求参数错误" {
		t.Errorf("message = %v, want '请求参数错误'", resp["message"])
	}
	if resp["data"] != nil {
		t.Error("data 应为 nil")
	}
}

func TestPaginated(t *testing.T) {
	c, w := setupTestContext()

	list := []string{"item1", "item2", "item3"}
	var total int64 = 10
	page := 1
	size := 3

	Paginated(c, list, total, page, size)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(10) {
		t.Errorf("total = %v, want 10", data["total"])
	}
	if data["page"] != float64(1) {
		t.Errorf("page = %v, want 1", data["page"])
	}
	if data["size"] != float64(3) {
		t.Errorf("size = %v, want 3", data["size"])
	}
	// totalPage = 10 / 3 = 3 余 1, 所以是 4
	if data["total_page"] != float64(4) {
		t.Errorf("total_page = %v, want 4", data["total_page"])
	}
}

func TestPaginatedExactDivision(t *testing.T) {
	c, w := setupTestContext()

	list := []string{"item1", "item2"}
	var total int64 = 6
	page := 1
	size := 3

	Paginated(c, list, total, page, size)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	// totalPage = 6 / 3 = 2, 余数为 0, 所以是 2
	if data["total_page"] != float64(2) {
		t.Errorf("total_page = %v, want 2", data["total_page"])
	}
}

func TestPaginatedEmptyList(t *testing.T) {
	c, w := setupTestContext()

	var total int64 = 0
	page := 1
	size := 10

	Paginated(c, nil, total, page, size)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data := resp["data"].(map[string]interface{})
	if data["total"] != float64(0) {
		t.Errorf("total = %v, want 0", data["total"])
	}
	if data["total_page"] != float64(0) {
		t.Errorf("total_page = %v, want 0", data["total_page"])
	}
}
