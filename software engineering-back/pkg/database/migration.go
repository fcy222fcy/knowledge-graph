package database

import (
	"log"
	"software_engineering/internal/model/entity"
)

// AutoMigrate 自动执行数据库迁移，创建或更新所有实体表
func AutoMigrate() {
	err := DB.AutoMigrate(
		&entity.User{},     // 学生表
		&entity.Teacher{},  // 教师表（相当于管理员）
		&entity.Document{},
		&entity.KnowledgePoint{},
		&entity.KnowledgeRelation{},
		&entity.KnowledgeBuild{},
		&entity.Question{},
		&entity.Quiz{},
		&entity.AskSession{},
		&entity.AskMessage{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	log.Println("database migration completed")
}
