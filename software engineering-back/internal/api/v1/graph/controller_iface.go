package graph

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// GraphService 定义图谱服务接口
type GraphService interface {
	GetGraphData(documentID uint, keyword, relationType string) (*response.GraphDataResponse, error) // 获取图谱数据
	BuildGraph(documentIDs []uint) (*response.BuildGraphResponse, error)                             // 构建图谱
	GetLatestBuildResult() (*response.BuildGraphResponse, error)                                     // 获取最新构建结果
	ListBuildHistory(page, size int) (*response.BuildHistoryResponse, error)                         // 列表构建历史
}

// GraphController 图谱控制器
type GraphController struct {
	graphService GraphService // 图谱服务
}

// NewGraphController 创建图谱控制器实例
func NewGraphController(graphService GraphService) *GraphController {
	return &GraphController{graphService: graphService}
}

// GetGraph 获取知识图谱数据
func (ctrl *GraphController) GetGraph(c *gin.Context) {
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	keyword := c.Query("keyword")
	relationType := c.Query("relation_type")
	resp, err := ctrl.graphService.GetGraphData(uint(documentID), keyword, relationType)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// BuildGraph 构建知识图谱
func (ctrl *GraphController) BuildGraph(c *gin.Context) {
	var req request.BuildGraphRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := ctrl.graphService.BuildGraph(req.DocumentIDs)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetLatestBuild 获取最近一次构建结果
func (ctrl *GraphController) GetLatestBuild(c *gin.Context) {
	resp, err := ctrl.graphService.GetLatestBuildResult()
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// ListBuildHistory 获取构建历史记录
func (ctrl *GraphController) ListBuildHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := ctrl.graphService.ListBuildHistory(page, size)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}
