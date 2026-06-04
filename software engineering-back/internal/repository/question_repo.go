package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateQuestion(q *model.Question) error {
	return database.DB.Create(q).Error
}

func FindQuestionByID(id uint) (*model.Question, error) {
	var q model.Question
	err := database.DB.First(&q, id).Error
	return &q, err
}

func UpdateQuestion(q *model.Question) error {
	return database.DB.Save(q).Error
}

func DeleteQuestion(id uint) error {
	return database.DB.Delete(&model.Question{}, id).Error
}

func ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64
	query := database.DB.Model(&model.Question{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if knowledgePointID > 0 {
		query = query.Where("knowledge_point_id = ?", knowledgePointID)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&questions).Error
	return questions, total, err
}
