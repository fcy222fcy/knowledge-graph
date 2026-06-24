package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// ListDocuments 获取文档列表（所有用户的）
func ListDocuments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	list, total, err := service.ListDocumentsAdmin(page, size, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// GetDocument 获取文档详情
func GetDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	doc, err := service.GetDocument(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "文档不存在")
		return
	}
	response.Success(c, doc)
}

// DeleteDocument 删除文档（管理员可以删除任何文档）
func DeleteDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := service.DeleteDocumentAdmin(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// ReviewDocument 审核文档（教师/管理员审核学生上传的文档）
func ReviewDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req request.ReviewDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := service.ReviewDocument(uint(id), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}
