package repository

import (
	"log"

	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

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

func FindKnowledgePointByID(id uint) (*entity.KnowledgePoint, error) {
	var kp entity.KnowledgePoint
	err := database.DB.First(&kp, id).Error
	return &kp, err
}

func UpdateKnowledgePoint(kp *entity.KnowledgePoint) error {
	return database.DB.Save(kp).Error
}

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

func FindRelationByID(id uint) (*entity.KnowledgeRelation, error) {
	var rel entity.KnowledgeRelation
	err := database.DB.First(&rel, id).Error
	return &rel, err
}

func UpdateRelation(rel *entity.KnowledgeRelation) error {
	return database.DB.Save(rel).Error
}

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

func GetAllKnowledgePoints() ([]entity.KnowledgePoint, error) {
	var points []entity.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}
