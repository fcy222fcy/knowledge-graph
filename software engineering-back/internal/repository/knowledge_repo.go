package repository

import (
	"log"

	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CreateKnowledgePoint 创建知识点，同时双写到 Neo4j（尽力而为）
func CreateKnowledgePoint(kp *entity.KnowledgePoint) error {
	if err := database.DB.Create(kp).Error; err != nil {
		return err
	}
	// Dual-write: sync to Neo4j (best effort)
	if err := CreateKnowledgePointInNeo4j(kp); err != nil {
		log.Printf("warning: neo4j dual-write failed for knowledge point %d: %v", kp.ID, err)
	}
	return nil
}

// FindKnowledgePointByID 根据 ID 查找知识点
func FindKnowledgePointByID(id uint) (*entity.KnowledgePoint, error) {
	var kp entity.KnowledgePoint
	err := database.DB.First(&kp, id).Error
	return &kp, err
}

// UpdateKnowledgePoint 更新知识点信息
func UpdateKnowledgePoint(kp *entity.KnowledgePoint) error {
	return database.DB.Save(kp).Error
}

// DeleteKnowledgePoint 删除知识点，同时从 Neo4j 同步删除（尽力而为）
func DeleteKnowledgePoint(id uint) error {
	if err := database.DB.Delete(&entity.KnowledgePoint{}, id).Error; err != nil {
		return err
	}
	// Dual-write: sync delete to Neo4j (best effort)
	if err := DeleteKnowledgePointFromNeo4j(id); err != nil {
		log.Printf("warning: neo4j dual-delete failed for knowledge point %d: %v", id, err)
	}
	return nil
}

// ListKnowledgePoints 分页查询知识点列表，支持按名称和文档 ID 过滤
func ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]entity.KnowledgePoint, int64, error) {
	var points []entity.KnowledgePoint
	var total int64
	query := database.DB.Model(&entity.KnowledgePoint{})
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

// CreateRelation 创建知识点关系，同时双写到 Neo4j（尽力而为）
func CreateRelation(rel *entity.KnowledgeRelation) error {
	if err := database.DB.Create(rel).Error; err != nil {
		return err
	}
	// Dual-write: sync to Neo4j (best effort)
	if err := CreateRelationInNeo4j(rel); err != nil {
		log.Printf("warning: neo4j dual-write failed for relation %d: %v", rel.ID, err)
	}
	return nil
}

// FindRelationByID 根据 ID 查找知识点关系
func FindRelationByID(id uint) (*entity.KnowledgeRelation, error) {
	var rel entity.KnowledgeRelation
	err := database.DB.First(&rel, id).Error
	return &rel, err
}

// UpdateRelation 更新知识点关系
func UpdateRelation(rel *entity.KnowledgeRelation) error {
	return database.DB.Save(rel).Error
}

// DeleteRelation 删除知识点关系，同时从 Neo4j 同步删除（尽力而为）
func DeleteRelation(id uint) error {
	if err := database.DB.Delete(&entity.KnowledgeRelation{}, id).Error; err != nil {
		return err
	}
	// Dual-write: sync delete to Neo4j (best effort)
	if err := DeleteRelationFromNeo4j(id); err != nil {
		log.Printf("warning: neo4j dual-delete failed for relation %d: %v", id, err)
	}
	return nil
}

// ListRelations 分页查询知识点关系列表，支持按知识点 ID 过滤
func ListRelations(page, size int, pointID uint) ([]entity.KnowledgeRelation, int64, error) {
	var rels []entity.KnowledgeRelation
	var total int64
	query := database.DB.Model(&entity.KnowledgeRelation{})
	if pointID > 0 {
		query = query.Where("source_id = ? OR target_id = ?", pointID, pointID)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&rels).Error
	return rels, total, err
}

// GetAllKnowledgePoints 获取所有知识点列表
func GetAllKnowledgePoints() ([]entity.KnowledgePoint, error) {
	var points []entity.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}
