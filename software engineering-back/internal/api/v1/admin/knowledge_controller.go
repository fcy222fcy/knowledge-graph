package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/repository"
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

// GetKnowledgeGraph 获取知识图谱数据
func GetKnowledgeGraph(c *gin.Context) {
	points, rels, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	nodes := make([]gin.H, len(points))
	for i, p := range points {
		nodes[i] = gin.H{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"document_id": p.DocumentID,
			"category":    p.Category,
		}
	}

	edges := make([]gin.H, len(rels))
	for i, r := range rels {
		edges[i] = gin.H{
			"id":            r.ID,
			"source":        r.SourceID,
			"target":        r.TargetID,
			"relation_type": r.RelationType,
			"description":   r.Description,
		}
	}

	response.Success(c, gin.H{
		"nodes": nodes,
		"edges": edges,
		"summary": gin.H{
			"node_count": len(nodes),
			"edge_count": len(edges),
		},
	})
}

// RebuildKnowledgeGraph 重建知识图谱并返回变化统计
func RebuildKnowledgeGraph(c *gin.Context) {
	var req struct {
		DocumentIDs []uint `json:"document_ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 记录旧数据数量
	oldPoints, oldRels, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取旧图谱数据失败: "+err.Error())
		return
	}
	oldPointCount := len(oldPoints)
	oldRelCount := len(oldRels)

	// 获取抽取服务
	extractionSvc := service.GetExtractionService()
	if extractionSvc == nil {
		response.Error(c, http.StatusInternalServerError, "抽取服务未初始化")
		return
	}

	// 逐文档重新抽取并构建图谱
	for _, docID := range req.DocumentIDs {
		doc, err := repository.FindDocumentByID(docID)
		if err != nil {
			continue
		}

		// 调用 extraction service 重新抽取
		result, err := extractionSvc.ExtractAndStoreToNeo4j(doc.Content, docID)
		if err != nil {
			continue
		}
		_ = result // 抽取结果已通过 ExtractAndStoreToNeo4j 存储

		// 调用 BuildGraph 构建图谱
		service.BuildGraph([]uint{docID})
	}

	// 获取新数据数量
	newPoints, newRels, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取新图谱数据失败: "+err.Error())
		return
	}
	newPointCount := len(newPoints)
	newRelCount := len(newRels)

	// 计算变化量
	added := newPointCount - oldPointCount
	if added < 0 {
		added = 0
	}
	deleted := oldPointCount - newPointCount
	if deleted < 0 {
		deleted = 0
	}
	modified := oldPointCount + newPointCount - added - deleted
	if modified < 0 {
		modified = 0
	}

	response.Success(c, gin.H{
		"old_point_count":   oldPointCount,
		"new_point_count":   newPointCount,
		"old_relation_count": oldRelCount,
		"new_relation_count": newRelCount,
		"added":             added,
		"deleted":           deleted,
		"modified":          modified,
	})
}
