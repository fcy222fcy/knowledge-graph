package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model/entity"
)

func CreateQuiz(quiz *entity.Quiz) error {
	return database.DB.Create(quiz).Error
}

func FindQuizByID(id uint) (*entity.Quiz, error) {
	var quiz entity.Quiz
	err := database.DB.First(&quiz, id).Error
	return &quiz, err
}

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

func CountCorrectByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&count).Error
	return count, err
}

func CountTotalByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
