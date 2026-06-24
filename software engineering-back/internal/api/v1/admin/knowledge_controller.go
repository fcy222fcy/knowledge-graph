package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// ListKnowledgePoints 获取知识点列表
func ListKnowledgePoints(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListKnowledgePointsSimple(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// DeleteKnowledgePoint 删除知识点
func DeleteKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := service.DeleteKnowledgePoint(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListKnowledgeRelations 获取知识点关系列表
func ListKnowledgeRelations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListKnowledgeRelationsSimple(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// DeleteKnowledgeRelation 删除知识点关系
func DeleteKnowledgeRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := service.DeleteRelation(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}
