package repository

import (
	"software_engineering/internal/model/entity"

	"gorm.io/gorm"
)

// QuestionRepository 定义题目仓库接口
type QuestionRepository interface {
	Create(q *entity.Question) error                // 创建题目
	FindByID(id uint) (*entity.Question, error)      // 根据 ID 查找题目
	Update(q *entity.Question) error                 // 更新题目
	Delete(id uint) error                            // 删除题目
	List(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]entity.Question, int64, error) // 分页查询
}

// QuizRepository 定义答题记录仓库接口
type QuizRepository interface {
	Create(quiz *entity.Quiz) error                                    // 创建答题记录
	FindByID(id uint) (*entity.Quiz, error)                           // 根据 ID 查找答题记录
	ListByUser(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]entity.Quiz, int64, error) // 分页获取用户答题记录
}

// MockQuestionRepository 题目仓库的 Mock 实现
type MockQuestionRepository struct {
	Questions map[uint]*entity.Question // 题目数据存储
	NextID    uint                       // 下一个自增 ID
	Err       error                      // 模拟错误
}

// NewMockQuestionRepository 创建 Mock 题目仓库
func NewMockQuestionRepository() *MockQuestionRepository {
	return &MockQuestionRepository{
		Questions: make(map[uint]*entity.Question),
		NextID:    1,
	}
}

func (m *MockQuestionRepository) Create(q *entity.Question) error {
	if m.Err != nil {
		return m.Err
	}
	q.ID = m.NextID
	m.NextID++
	m.Questions[q.ID] = q
	return nil
}

func (m *MockQuestionRepository) FindByID(id uint) (*entity.Question, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	q, ok := m.Questions[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return q, nil
}

func (m *MockQuestionRepository) Update(q *entity.Question) error {
	if m.Err != nil {
		return m.Err
	}
	m.Questions[q.ID] = q
	return nil
}

func (m *MockQuestionRepository) Delete(id uint) error {
	if m.Err != nil {
		return m.Err
	}
	delete(m.Questions, id)
	return nil
}

func (m *MockQuestionRepository) List(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]entity.Question, int64, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	var questions []entity.Question
	for _, q := range m.Questions {
		questions = append(questions, *q)
	}
	return questions, int64(len(questions)), nil
}

// AddQuestion 辅助方法：添加测试题目
func (m *MockQuestionRepository) AddQuestion(id uint, title, answer, questionType, difficulty, options, explanation string) {
	q := &entity.Question{
		Title:       title,
		Type:        questionType,
		Difficulty:  difficulty,
		Options:     options,
		Answer:      answer,
		Explanation: explanation,
	}
	q.ID = id
	m.Questions[id] = q
	if id >= m.NextID {
		m.NextID = id + 1
	}
}

// SetError 辅助方法：设置错误
func (m *MockQuestionRepository) SetError(err error) {
	m.Err = err
}

// ClearError 辅助方法：清除错误
func (m *MockQuestionRepository) ClearError() {
	m.Err = nil
}

// MockQuizRepository 答题记录仓库的 Mock 实现
type MockQuizRepository struct {
	Quizzes map[uint]*entity.Quiz // 答题记录数据存储
	NextID  uint                   // 下一个自增 ID
	Err     error                  // 模拟错误
}

// NewMockQuizRepository 创建 Mock 答题记录仓库
func NewMockQuizRepository() *MockQuizRepository {
	return &MockQuizRepository{
		Quizzes: make(map[uint]*entity.Quiz),
		NextID:  1,
	}
}

func (m *MockQuizRepository) Create(quiz *entity.Quiz) error {
	if m.Err != nil {
		return m.Err
	}
	quiz.ID = m.NextID
	m.NextID++
	m.Quizzes[quiz.ID] = quiz
	return nil
}

func (m *MockQuizRepository) FindByID(id uint) (*entity.Quiz, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	quiz, ok := m.Quizzes[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return quiz, nil
}

func (m *MockQuizRepository) ListByUser(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]entity.Quiz, int64, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	var quizzes []entity.Quiz
	for _, quiz := range m.Quizzes {
		if quiz.UserID == userID {
			if isCorrect != nil && quiz.IsCorrect != *isCorrect {
				continue
			}
			quizzes = append(quizzes, *quiz)
		}
	}
	return quizzes, int64(len(quizzes)), nil
}

// AddQuiz 辅助方法：添加测试答题记录
func (m *MockQuizRepository) AddQuiz(id, questionID, userID uint, userAnswer string, isCorrect bool) {
	quiz := &entity.Quiz{
		QuestionID: questionID,
		UserID:     userID,
		UserAnswer: userAnswer,
		IsCorrect:  isCorrect,
	}
	quiz.ID = id
	m.Quizzes[id] = quiz
	if id >= m.NextID {
		m.NextID = id + 1
	}
}

// SetError 辅助方法：设置错误
func (m *MockQuizRepository) SetError(err error) {
	m.Err = err
}

// ClearError 辅助方法：清除错误
func (m *MockQuizRepository) ClearError() {
	m.Err = nil
}
