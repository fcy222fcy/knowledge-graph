package repository

import (
	"log"

	"software_engineering/internal/model/entity"
	"software_engineering/pkg/database"
	"gorm.io/gorm"
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

// CreateKnowledgePointWithRelations 在一个事务中批量创建知识点及其关系
// 示例：展示如何使用 GORM 事务保证原子性
func CreateKnowledgePointWithRelations(kp *entity.KnowledgePoint, rels []entity.KnowledgeRelation) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建知识点
		if err := tx.Create(kp).Error; err != nil {
			return err
		}
		// 2. 批量创建关系（填充 source_id）
		for i := range rels {
			rels[i].SourceID = kp.ID
		}
		if len(rels) > 0 {
			if err := tx.Create(&rels).Error; err != nil {
				return err // 回滚知识点和已创建的关系
			}
		}
		return nil // 提交事务
	})
}

// FindKnowledgePointByID 根据 ID 查找知识点
func FindKnowledgePointByID(id uint) (*entity.KnowledgePoint, error) {
	var kp entity.KnowledgePoint
	err := database.DB.First(&kp, id).Error
	return &kp, err
}

// FindKnowledgePointByName 根据名称和文档 ID 查找知识点
func FindKnowledgePointByName(name string, documentID uint) (*entity.KnowledgePoint, error) {
	var kp entity.KnowledgePoint
	err := database.DB.Where("name = ? AND document_id = ?", name, documentID).First(&kp).Error
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

// CountKnowledgePoints 统计知识点总数
func CountKnowledgePoints() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.KnowledgePoint{}).Count(&count).Error
	return count, err
}

// CountRelations 统计知识点关系总数
func CountRelations() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.KnowledgeRelation{}).Count(&count).Error
	return count, err
}

// ListKnowledgePointsAdmin 管理员获取知识点列表
func ListKnowledgePointsAdmin(page, size int) ([]entity.KnowledgePoint, int64, error) {
	var points []entity.KnowledgePoint
	var total int64
	database.DB.Model(&entity.KnowledgePoint{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&points).Error
	return points, total, err
}

// ListRelationsAdmin 管理员获取关系列表
func ListRelationsAdmin(page, size int) ([]entity.KnowledgeRelation, int64, error) {
	var rels []entity.KnowledgeRelation
	var total int64
	database.DB.Model(&entity.KnowledgeRelation{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&rels).Error
	return rels, total, err
}
