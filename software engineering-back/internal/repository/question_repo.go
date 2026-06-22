package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CreateQuestion 创建题目
func CreateQuestion(q *entity.Question) error {
	return database.DB.Create(q).Error
}

// FindQuestionByID 根据 ID 查找题目
func FindQuestionByID(id uint) (*entity.Question, error) {
	var q entity.Question
	err := database.DB.First(&q, id).Error
	return &q, err
}

// UpdateQuestion 更新题目信息
func UpdateQuestion(q *entity.Question) error {
	return database.DB.Save(q).Error
}

// DeleteQuestion 删除题目
func DeleteQuestion(id uint) error {
	return database.DB.Delete(&entity.Question{}, id).Error
}

// ListQuestions 分页查询题目列表，支持按标题、知识点和难度过滤
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
