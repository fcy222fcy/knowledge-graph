package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateKnowledgeBuild(build *model.KnowledgeBuild) error {
	return database.DB.Create(build).Error
}

func GetLatestBuild() (*model.KnowledgeBuild, error) {
	var build model.KnowledgeBuild
	err := database.DB.Order("created_at DESC").First(&build).Error
	return &build, err
}

func ListBuilds(page, size int) ([]model.KnowledgeBuild, int64, error) {
	var builds []model.KnowledgeBuild
	var total int64
	database.DB.Model(&model.KnowledgeBuild{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&builds).Error
	return builds, total, err
}

func GetAllKnowledgePointsForGraph() ([]model.KnowledgePoint, error) {
	var points []model.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}

func GetAllRelationsForGraph() ([]model.KnowledgeRelation, error) {
	var rels []model.KnowledgeRelation
	err := database.DB.Find(&rels).Error
	return rels, err
}

func FindKnowledgePointsByIDs(ids []uint) ([]model.KnowledgePoint, error) {
	var points []model.KnowledgePoint
	err := database.DB.Where("id IN ?", ids).Find(&points).Error
	return points, err
}
