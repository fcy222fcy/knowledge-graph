package quiz

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// QuizService 定义答题服务接口
type QuizService interface {
	SubmitQuiz(userID uint, req request.SubmitQuizRequest) (*response.QuizResponse, error)                                          // 提交答题
	GetQuizDetail(id uint) (*response.QuizResponse, error)                                                                         // 获取答题详情
	ListQuizHistory(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]response.QuizResponse, int64, error)   // 列表答题历史
}

// QuizController 答题控制器
type QuizController struct {
	quizService QuizService // 答题服务
}

// NewQuizController 创建答题控制器实例
func NewQuizController(quizService QuizService) *QuizController {
	return &QuizController{quizService: quizService}
}

// SubmitQuiz 提交答题
func (ctrl *QuizController) SubmitQuiz(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := ctrl.quizService.SubmitQuiz(userID, req)
	if err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// GetQuizDetail 获取答题记录详情
func (ctrl *QuizController) GetQuizDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := ctrl.quizService.GetQuizDetail(uint(id))
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// ListQuizHistory 获取答题历史记录
func (ctrl *QuizController) ListQuizHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	list, total, err := ctrl.quizService.ListQuizHistory(userID, page, size, uint(knowledgePointID), nil)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Paginated(c, list, total, page, size)
}
