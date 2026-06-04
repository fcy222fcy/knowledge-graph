package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateKnowledgePoint(kp *model.KnowledgePoint) error {
	return database.DB.Create(kp).Error
}

func FindKnowledgePointByID(id uint) (*model.KnowledgePoint, error) {
	var kp model.KnowledgePoint
	err := database.DB.First(&kp, id).Error
	return &kp, err
}

func UpdateKnowledgePoint(kp *model.KnowledgePoint) error {
	return database.DB.Save(kp).Error
}

func DeleteKnowledgePoint(id uint) error {
	return database.DB.Delete(&model.KnowledgePoint{}, id).Error
}

func ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]model.KnowledgePoint, int64, error) {
	var points []model.KnowledgePoint
	var total int64
	query := database.DB.Model(&model.KnowledgePoint{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if documentID > 0 {
		query = query.Where("document_id = ?", documentID)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&points).Error
	return points, total, err
}

func CreateRelation(rel *model.KnowledgeRelation) error {
	return database.DB.Create(rel).Error
}

func FindRelationByID(id uint) (*model.KnowledgeRelation, error) {
	var rel model.KnowledgeRelation
	err := database.DB.First(&rel, id).Error
	return &rel, err
}

func UpdateRelation(rel *model.KnowledgeRelation) error {
	return database.DB.Save(rel).Error
}

func DeleteRelation(id uint) error {
	return database.DB.Delete(&model.KnowledgeRelation{}, id).Error
}

func ListRelations(page, size int, pointID uint) ([]model.KnowledgeRelation, int64, error) {
	var rels []model.KnowledgeRelation
	var total int64
	query := database.DB.Model(&model.KnowledgeRelation{})
	if pointID > 0 {
		query = query.Where("source_id = ? OR target_id = ?", pointID, pointID)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&rels).Error
	return rels, total, err
}

func GetAllKnowledgePoints() ([]model.KnowledgePoint, error) {
	var points []model.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}