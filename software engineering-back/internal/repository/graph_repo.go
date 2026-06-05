package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"software_engineering/pkg/database"
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

// --- Neo4j Cypher queries ---

func GetAllGraphDataFromNeo4j() ([]entity.KnowledgePoint, []entity.KnowledgeRelation, error) {
	if !database.IsNeo4jAvailable() {
		points, err := GetAllKnowledgePointsForGraph()
		if err != nil {
			return nil, nil, err
		}
		rels, err := GetAllRelationsForGraph()
		return points, rels, err
	}

	session := database.Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(context.Background())

	ctx := context.Background()

	// Fetch all KnowledgePoint nodes
	pointsResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `MATCH (n:KnowledgePoint) RETURN n.id AS id, n.name AS name, n.description AS description, n.document_id AS document_id, n.category AS category`, nil)
		if err != nil {
			return nil, err
		}
		var points []entity.KnowledgePoint
		for result.Next(ctx) {
			record := recordToMap(result.Record())
			points = append(points, entity.KnowledgePoint{
				BaseModel: entity.BaseModel{
					ID: toUint(record["id"]),
				},
				Name:        toString(record["name"]),
				Description: toString(record["description"]),
				DocumentID:  toUint(record["document_id"]),
				Category:    toString(record["category"]),
			})
		}
		return points, result.Err()
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read nodes from neo4j: %w", err)
	}

	// Fetch all KNOWS relationships (using a generic relationship type)
	relsResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `MATCH (a:KnowledgePoint)-[r:RELATED|DEPENDS_ON|PART_OF]->(b:KnowledgePoint) RETURN r.id AS id, a.id AS source_id, b.id AS target_id, r.type AS relation_type, r.description AS description`, nil)
		if err != nil {
			return nil, err
		}
		var rels []entity.KnowledgeRelation
		for result.Next(ctx) {
			record := recordToMap(result.Record())
			rels = append(rels, entity.KnowledgeRelation{
				BaseModel: entity.BaseModel{
					ID: toUint(record["id"]),
				},
				SourceID:     toUint(record["source_id"]),
				TargetID:     toUint(record["target_id"]),
				RelationType: toString(record["relation_type"]),
				Description:  toString(record["description"]),
			})
		}
		return rels, result.Err()
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read relations from neo4j: %w", err)
	}

	points, _ := pointsResult.([]entity.KnowledgePoint)
	rels, _ := relsResult.([]entity.KnowledgeRelation)
	return points, rels, nil
}

func CreateKnowledgePointInNeo4j(kp *entity.KnowledgePoint) error {
	if !database.IsNeo4jAvailable() {
		return nil
	}

	session := database.Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(context.Background())

	ctx := context.Background()
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `CREATE (n:KnowledgePoint {id: $id, name: $name, description: $description, document_id: $document_id, category: $category})`,
			map[string]any{
				"id":          int64(kp.ID),
				"name":        kp.Name,
				"description": kp.Description,
				"document_id": int64(kp.DocumentID),
				"category":    kp.Category,
			})
		return nil, err
	})
	if err != nil {
		log.Printf("warning: failed to create knowledge point in neo4j: %v", err)
	}
	return err
}

func CreateRelationInNeo4j(rel *entity.KnowledgeRelation) error {
	if !database.IsNeo4jAvailable() {
		return nil
	}

	// Validate relation type to prevent Cypher injection
	validRelationTypes := map[string]bool{
		"RELATED":    true,
		"DEPENDS_ON": true,
		"PART_OF":    true,
		"IMPLEMENTS": true,
		"EXTENDS":    true,
		"USES":       true,
	}
	if !validRelationTypes[rel.RelationType] {
		return fmt.Errorf("invalid relation type: %s", rel.RelationType)
	}

	session := database.Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(context.Background())

	ctx := context.Background()
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Use the relation type as the Neo4j relationship type
		cypher := fmt.Sprintf(`MATCH (a:KnowledgePoint {id: $source_id}), (b:KnowledgePoint {id: $target_id}) CREATE (a)-[r:%s {id: $id, description: $description}]->(b)`, rel.RelationType)
		_, err := tx.Run(ctx, cypher, map[string]any{
			"id":          int64(rel.ID),
			"source_id":   int64(rel.SourceID),
			"target_id":   int64(rel.TargetID),
			"description": rel.Description,
		})
		return nil, err
	})
	if err != nil {
		log.Printf("warning: failed to create relation in neo4j: %v", err)
	}
	return err
}

func DeleteKnowledgePointFromNeo4j(id uint) error {
	if !database.IsNeo4jAvailable() {
		return nil
	}

	session := database.Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(context.Background())

	ctx := context.Background()
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `MATCH (n:KnowledgePoint {id: $id}) DETACH DELETE n`, map[string]any{
			"id": int64(id),
		})
		return nil, err
	})
	if err != nil {
		log.Printf("warning: failed to delete knowledge point from neo4j: %v", err)
	}
	return err
}

func DeleteRelationFromNeo4j(id uint) error {
	if !database.IsNeo4jAvailable() {
		return nil
	}

	session := database.Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(context.Background())

	ctx := context.Background()
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `MATCH ()-[r]->() WHERE r.id = $id DELETE r`, map[string]any{
			"id": int64(id),
		})
		return nil, err
	})
	if err != nil {
		log.Printf("warning: failed to delete relation from neo4j: %v", err)
	}
	return err
}

// --- helper functions ---

func recordToMap(record *neo4j.Record) map[string]any {
	m := make(map[string]any)
	for _, key := range record.Keys {
		val, _ := record.Get(key)
		m[key] = val
	}
	return m
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func toUint(v any) uint {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case int64:
		return uint(n)
	case int:
		return uint(n)
	case float64:
		return uint(n)
	default:
		return 0
	}
}
