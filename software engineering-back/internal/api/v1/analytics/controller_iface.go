package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// AnalyticsService 定义分析服务接口
type AnalyticsService interface {
	GetOverview(userID uint) (*response.OverviewResponse, error)                             // 获取学习概览
	GetHotKnowledgePoints(limit int) ([]response.HotKnowledgePoint, error)                  // 获取热门知识点
	GetKnowledgeMastery(userID uint) ([]response.KnowledgeMastery, error)                   // 获取知识点掌握程度
	GetWeakPoints(userID uint, limit int) ([]response.WeakPoint, error)                     // 获取薄弱知识点
	GetTrends(userID uint, days int) (*response.TrendData, error)                           // 获取学习趋势
}

// AnalyticsController 分析控制器
type AnalyticsController struct {
	analyticsService AnalyticsService // 分析服务
}

// NewAnalyticsController 创建分析控制器实例
func NewAnalyticsController(analyticsService AnalyticsService) *AnalyticsController {
	return &AnalyticsController{analyticsService: analyticsService}
}

// GetOverview 获取学习概览数据
func (ctrl *AnalyticsController) GetOverview(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := ctrl.analyticsService.GetOverview(userID)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetHotKnowledgePoints 获取热门知识点列表
func (ctrl *AnalyticsController) GetHotKnowledgePoints(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := ctrl.analyticsService.GetHotKnowledgePoints(limit)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetKnowledgeMastery 获取知识点掌握程度
func (ctrl *AnalyticsController) GetKnowledgeMastery(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := ctrl.analyticsService.GetKnowledgeMastery(userID)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetWeakPoints 获取薄弱知识点
func (ctrl *AnalyticsController) GetWeakPoints(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := ctrl.analyticsService.GetWeakPoints(userID, limit)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetTrends 获取学习趋势数据
func (ctrl *AnalyticsController) GetTrends(c *gin.Context) {
	userID := c.GetUint("user_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	resp, err := ctrl.analyticsService.GetTrends(userID, days)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}
