package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model/entity"
)

func CreateQuestion(q *entity.Question) error {
	return database.DB.Create(q).Error
}

func FindQuestionByID(id uint) (*entity.Question, error) {
	var q entity.Question
	err := database.DB.First(&q, id).Error
	return &q, err
}

func UpdateQuestion(q *entity.Question) error {
	return database.DB.Save(q).Error
}

func DeleteQuestion(id uint) error {
	return database.DB.Delete(&entity.Question{}, id).Error
}

func ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]entity.Question, int64, error) {
	var questions []entity.Question
	var total int64
	query := database.DB.Model(&entity.Question{})
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
