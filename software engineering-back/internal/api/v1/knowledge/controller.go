package knowledge

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// CreateKnowledgePoint 创建知识点
func CreateKnowledgePoint(c *gin.Context) {
	var req request.CreateKnowledgePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateKnowledgePoint(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// GetKnowledgePoint 获取知识点详情
func GetKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetKnowledgePoint(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, resp)
}

// UpdateKnowledgePoint 更新知识点
func UpdateKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateKnowledgePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateKnowledgePoint(uint(id), req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteKnowledgePoint 删除知识点
func DeleteKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteKnowledgePoint(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListKnowledgePoints 获取知识点列表
func ListKnowledgePoints(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	list, total, err := service.ListKnowledgePoints(page, size, keyword, uint(documentID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// CreateRelation 创建知识点关系
func CreateRelation(c *gin.Context) {
	var req request.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateRelation(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

// UpdateRelation 更新知识点关系
func UpdateRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateRelation(uint(id), req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteRelation 删除知识点关系
func DeleteRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteRelation(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListRelations 获取知识点关系列表
func ListRelations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	pointID, _ := strconv.Atoi(c.Query("point_id"))
	list, total, err := service.ListRelations(page, size, uint(pointID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}
