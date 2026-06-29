package database

import (
	"fmt"
	"log"
	"software_engineering/internal/model/entity"
)

// indexDef 描述需要创建的复合索引
type indexDef struct {
	table string
	name  string
	cols  string // 列名列表，如 "user_id, is_correct"
}

// requiredIndexes 定义所有需要的复合索引（单列索引由 AutoMigrate 的 gorm:"index" tag 处理）
var requiredIndexes = []indexDef{
	{"quizzes", "idx_quizzes_user_id_is_correct", "user_id, is_correct"},
	{"quizzes", "idx_quizzes_user_id_created_at", "user_id, created_at"},
	{"ask_messages", "idx_ask_messages_session_id_role", "session_id, role"},
	{"ask_messages", "idx_ask_messages_session_id_created_at", "session_id, created_at"},
	{"documents", "idx_documents_user_id_status", "user_id, status"},
}

// AutoMigrate 自动执行数据库迁移，创建或更新所有实体表，并确保复合索引存在
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
		&entity.Assignment{},
		&entity.AssignmentQuestion{},
		&entity.AssignmentSubmission{},
		&entity.AskSession{},
		&entity.AskMessage{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// 创建复合索引（幂等：CREATE INDEX IF NOT EXISTS）
	if err := ensureIndexes(); err != nil {
		log.Fatalf("failed to create indexes: %v", err)
	}

	log.Println("database migration completed")
}

// ensureIndexes 安全地创建复合索引（兼容 MySQL 5.7+）
func ensureIndexes() error {
	for _, idx := range requiredIndexes {
		// 先检查索引是否存在
		var count int
		checkSQL := fmt.Sprintf(
			"SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = '%s' AND index_name = '%s'",
			idx.table, idx.name,
		)
		if err := DB.Raw(checkSQL).Scan(&count).Error; err != nil {
			log.Printf("warning: failed to check index %s: %v", idx.name, err)
			continue
		}

		// 如果索引不存在，则创建
		if count == 0 {
			createSQL := fmt.Sprintf(
				"CREATE INDEX %s ON %s (%s)",
				idx.name, idx.table, idx.cols,
			)
			if err := DB.Exec(createSQL).Error; err != nil {
				log.Printf("warning: failed to create index %s: %v", idx.name, err)
				continue
			}
			log.Printf("created index %s on %s(%s)", idx.name, idx.table, idx.cols)
		}
	}
	return nil
}
