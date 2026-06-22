package knowledge

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// KnowledgeService 定义知识点服务接口
type KnowledgeService interface {
	CreateKnowledgePoint(req request.CreateKnowledgePointRequest) (uint, error)                                                    // 创建知识点
	GetKnowledgePoint(id uint) (*response.KnowledgePointResponse, error)                                                           // 获取知识点
	UpdateKnowledgePoint(id uint, req request.UpdateKnowledgePointRequest) error                                                    // 更新知识点
	DeleteKnowledgePoint(id uint) error                                                                                             // 删除知识点
	ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]response.KnowledgePointResponse, int64, error)          // 列表知识点
	CreateRelation(req request.CreateRelationRequest) (uint, error)                                                                // 创建关系
	UpdateRelation(id uint, req request.UpdateRelationRequest) error                                                                // 更新关系
	DeleteRelation(id uint) error                                                                                                   // 删除关系
	ListRelations(page, size int, pointID uint) ([]response.KnowledgeRelationResponse, int64, error)                               // 列表关系
}

// KnowledgeController 知识点控制器
type KnowledgeController struct {
	knowledgeService KnowledgeService // 知识点服务
}

// NewKnowledgeController 创建知识点控制器实例
func NewKnowledgeController(knowledgeService KnowledgeService) *KnowledgeController {
	return &KnowledgeController{knowledgeService: knowledgeService}
}

// CreateKnowledgePoint 创建知识点
func (ctrl *KnowledgeController) CreateKnowledgePoint(c *gin.Context) {
	var req request.CreateKnowledgePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := ctrl.knowledgeService.CreateKnowledgePoint(req)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, gin.H{"id": id})
}

// GetKnowledgePoint 获取知识点详情
func (ctrl *KnowledgeController) GetKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := ctrl.knowledgeService.GetKnowledgePoint(uint(id))
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// UpdateKnowledgePoint 更新知识点
func (ctrl *KnowledgeController) UpdateKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateKnowledgePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := ctrl.knowledgeService.UpdateKnowledgePoint(uint(id), req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// DeleteKnowledgePoint 删除知识点
func (ctrl *KnowledgeController) DeleteKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.knowledgeService.DeleteKnowledgePoint(uint(id)); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// ListKnowledgePoints 获取知识点列表
func (ctrl *KnowledgeController) ListKnowledgePoints(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	list, total, err := ctrl.knowledgeService.ListKnowledgePoints(page, size, keyword, uint(documentID))
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}

// CreateRelation 创建知识点关系
func (ctrl *KnowledgeController) CreateRelation(c *gin.Context) {
	var req request.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := ctrl.knowledgeService.CreateRelation(req)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, gin.H{"id": id})
}

// UpdateRelation 更新知识点关系
func (ctrl *KnowledgeController) UpdateRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := ctrl.knowledgeService.UpdateRelation(uint(id), req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// DeleteRelation 删除知识点关系
func (ctrl *KnowledgeController) DeleteRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.knowledgeService.DeleteRelation(uint(id)); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// ListRelations 获取知识点关系列表
func (ctrl *KnowledgeController) ListRelations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	pointID, _ := strconv.Atoi(c.Query("point_id"))
	list, total, err := ctrl.knowledgeService.ListRelations(page, size, uint(pointID))
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}
