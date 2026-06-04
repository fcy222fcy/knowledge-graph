package database

import (
	"log"
	"software_engineering/internal/model"
)

func AutoMigrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Document{},
		&model.KnowledgePoint{},
		&model.KnowledgeRelation{},
		&model.KnowledgeBuild{},
		&model.Question{},
		&model.Quiz{},
		&model.AskSession{},
		&model.AskMessage{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	log.Println("database migration completed")
}
