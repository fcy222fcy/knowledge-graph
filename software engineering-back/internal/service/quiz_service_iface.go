package service

import (
	"errors"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

// QuizService 定义答题服务接口
type QuizService interface {
	SubmitQuiz(userID uint, req request.SubmitQuizRequest) (*response.QuizResponse, error)   // 提交答题
	GetQuizDetail(id uint) (*response.QuizResponse, error)                                  // 获取答题记录详情
	ListQuizHistory(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]response.QuizResponse, int64, error) // 分页获取答题历史
}

// QuizServiceImpl 答题服务实现
type QuizServiceImpl struct {
	questionRepo repository.QuestionRepository // 题目仓库
	quizRepo     repository.QuizRepository     // 答题仓库
}

// NewQuizService 创建答题服务实例
func NewQuizService(questionRepo repository.QuestionRepository, quizRepo repository.QuizRepository) *QuizServiceImpl {
	return &QuizServiceImpl{questionRepo: questionRepo, quizRepo: quizRepo}
}

// SubmitQuiz 提交答题，自动判断对错并保存记录
func (s *QuizServiceImpl) SubmitQuiz(userID uint, req request.SubmitQuizRequest) (*response.QuizResponse, error) {
	question, err := s.questionRepo.FindByID(req.QuestionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	isCorrect := question.Answer == req.UserAnswer

	quiz := &entity.Quiz{
		QuestionID: req.QuestionID,
		UserID:     userID,
		UserAnswer: req.UserAnswer,
		IsCorrect:  isCorrect,
	}
	if err := s.quizRepo.Create(quiz); err != nil {
		return nil, err
	}

	return &response.QuizResponse{
		QuizID:        quiz.ID,
		QuestionID:    question.ID,
		UserAnswer:    req.UserAnswer,
		CorrectAnswer: question.Answer,
		IsCorrect:     isCorrect,
		Explanation:   question.Explanation,
		CreatedAt:     quiz.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// GetQuizDetail 获取答题记录详情，包含题目信息和用户答案
func (s *QuizServiceImpl) GetQuizDetail(id uint) (*response.QuizResponse, error) {
	quiz, err := s.quizRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("答题记录不存在")
	}
	question, err := s.questionRepo.FindByID(quiz.QuestionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	return &response.QuizResponse{
		QuizID:        quiz.ID,
		QuestionID:    question.ID,
		QuestionTitle: question.Title,
		Type:          question.Type,
		Difficulty:    question.Difficulty,
		Options:       parseOptions(question.Options),
		UserAnswer:    quiz.UserAnswer,
		CorrectAnswer: question.Answer,
		IsCorrect:     quiz.IsCorrect,
		Explanation:   question.Explanation,
		CreatedAt:     quiz.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// ListQuizHistory 分页获取用户答题历史，支持按知识点和正确性过滤
func (s *QuizServiceImpl) ListQuizHistory(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]response.QuizResponse, int64, error) {
	quizzes, total, err := s.quizRepo.ListByUser(userID, page, size, knowledgePointID, isCorrect)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.QuizResponse, len(quizzes))
	for i, q := range quizzes {
		question, _ := s.questionRepo.FindByID(q.QuestionID)
		item := response.QuizResponse{
			QuizID:     q.ID,
			QuestionID: q.QuestionID,
			UserAnswer: q.UserAnswer,
			IsCorrect:  q.IsCorrect,
			CreatedAt:  q.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if question != nil {
			item.QuestionTitle = question.Title
			item.KnowledgePointID = question.KnowledgePointID
		}
		list[i] = item
	}
	return list, total, nil
}
