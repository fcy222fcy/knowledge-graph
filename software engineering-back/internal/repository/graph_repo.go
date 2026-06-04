package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model/entity"
)

func CreateKnowledgeBuild(build *entity.KnowledgeBuild) error {
	return database.DB.Create(build).Error
}

func GetLatestBuild() (*entity.KnowledgeBuild, error) {
	var build entity.KnowledgeBuild
	err := database.DB.Order("created_at DESC").First(&build).Error
	return &build, err
}

func ListBuilds(page, size int) ([]entity.KnowledgeBuild, int64, error) {
	var builds []entity.KnowledgeBuild
	var total int64
	database.DB.Model(&entity.KnowledgeBuild{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&builds).Error
	return builds, total, err
}

func GetAllKnowledgePointsForGraph() ([]entity.KnowledgePoint, error) {
	var points []entity.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}

func GetAllRelationsForGraph() ([]entity.KnowledgeRelation, error) {
	var rels []entity.KnowledgeRelation
	err := database.DB.Find(&rels).Error
	return rels, err
}

func FindKnowledgePointsByIDs(ids []uint) ([]entity.KnowledgePoint, error) {
	var points []entity.KnowledgePoint
	err := database.DB.Where("id IN ?", ids).Find(&points).Error
	return points, err
}
