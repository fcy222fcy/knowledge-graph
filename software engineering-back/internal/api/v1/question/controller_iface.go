package question

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// QuestionService 定义题目服务接口
type QuestionService interface {
	CreateQuestion(req request.CreateQuestionRequest) (uint, error)                                                               // 创建题目
	GetQuestion(id uint, includeAnswer bool) (*response.QuestionResponse, error)                                                  // 获取题目
	UpdateQuestion(id uint, req request.UpdateQuestionRequest) error                                                              // 更新题目
	DeleteQuestion(id uint) error                                                                                                 // 删除题目
	ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]response.QuestionResponse, int64, error) // 列表题目
}

// QuestionController 题目控制器
type QuestionController struct {
	questionService QuestionService // 题目服务
}

// NewQuestionController 创建题目控制器实例
func NewQuestionController(questionService QuestionService) *QuestionController {
	return &QuestionController{questionService: questionService}
}

// CreateQuestion 创建题目
func (ctrl *QuestionController) CreateQuestion(c *gin.Context) {
	var req request.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := ctrl.questionService.CreateQuestion(req)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, gin.H{"id": id})
}

// GetQuestion 获取题目详情
func (ctrl *QuestionController) GetQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := ctrl.questionService.GetQuestion(uint(id), true)
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// UpdateQuestion 更新题目
func (ctrl *QuestionController) UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := ctrl.questionService.UpdateQuestion(uint(id), req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// DeleteQuestion 删除题目
func (ctrl *QuestionController) DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.questionService.DeleteQuestion(uint(id)); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// ListQuestions 获取题目列表
func (ctrl *QuestionController) ListQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	difficulty := c.Query("difficulty")
	list, total, err := ctrl.questionService.ListQuestions(page, size, keyword, uint(knowledgePointID), difficulty)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}
