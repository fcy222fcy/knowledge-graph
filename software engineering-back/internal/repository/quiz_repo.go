package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CreateQuiz 创建答题记录
func CreateQuiz(quiz *entity.Quiz) error {
	return database.DB.Create(quiz).Error
}

// FindQuizByID 根据 ID 查找答题记录
func FindQuizByID(id uint) (*entity.Quiz, error) {
	var quiz entity.Quiz
	err := database.DB.First(&quiz, id).Error
	return &quiz, err
}

// ListQuizzesByUser 分页获取用户的答题记录，支持按知识点和正确性过滤
func ListQuizzesByUser(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]entity.Quiz, int64, error) {
	var quizzes []entity.Quiz
	var total int64
	query := database.DB.Model(&entity.Quiz{}).Where("user_id = ?", userID)
	if knowledgePointID > 0 {
		query = query.Joins("JOIN questions ON quizzes.question_id = questions.id AND questions.knowledge_point_id = ?", knowledgePointID)
	}
	if isCorrect != nil {
		query = query.Where("is_correct = ?", *isCorrect)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&quizzes).Error
	return quizzes, total, err
}

// CountCorrectByUser 统计用户答对的题目数量
func CountCorrectByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&count).Error
	return count, err
}

// CountTotalByUser 统计用户答题总数
func CountTotalByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountQuizzes 统计答题记录总数
func CountQuizzes() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Count(&count).Error
	return count, err
}
