package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel embeds gorm.Model for common fields.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ─── User ─────────────────────────────────────────────

type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1=active, 0=disabled
}

// ─── Document ─────────────────────────────────────────

type Document struct {
	BaseModel
	Title       string `gorm:"size:200;not null" json:"title"`
	Description string `gorm:"size:500" json:"description"`
	Filename    string `gorm:"size:200;not null" json:"filename"`
	FilePath    string `gorm:"size:500;not null" json:"-"`
	FileSize    int64  `json:"file_size"`
	FileType    string `gorm:"size:20" json:"file_type"`
	Content     string `gorm:"type:longtext" json:"-"`
	Status      string `gorm:"size:20;default:pending" json:"status"` // pending/processing/completed/failed
}

// ─── Knowledge ────────────────────────────────────────

type KnowledgePoint struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `gorm:"size:50" json:"category"`
}

type KnowledgeRelation struct {
	BaseModel
	SourceID     uint   `json:"source_id"`
	TargetID     uint   `json:"target_id"`
	RelationType string `gorm:"size:20;not null" json:"relation_type"` // RELATED/DEPENDS_ON/PART_OF
	Description  string `gorm:"size:500" json:"description"`
}

type KnowledgeBuild struct {
	BaseModel
	DocumentIDs    string `gorm:"size:500" json:"document_ids"` // comma-separated
	CreatedPoints  int    `json:"created_points"`
	CreatedRelations int  `json:"created_relations"`
	ChunkCount     int    `json:"chunk_count"`
	VectorCount    int    `json:"vector_count"`
	Status         string `gorm:"size:20;default:completed" json:"status"`
	Message        string `gorm:"size:500" json:"message"`
}

// ─── Question ─────────────────────────────────────────

type Question struct {
	BaseModel
	Title             string `gorm:"size:500;not null" json:"title"`
	Type              string `gorm:"size:20;not null" json:"type"` // single/multiple
	Difficulty        string `gorm:"size:20;not null" json:"difficulty"` // easy/medium/hard
	KnowledgePointID  uint   `json:"knowledge_point_id"`
	Options           string `gorm:"type:text" json:"-"` // JSON array stored as text
	Answer            string `gorm:"size:20;not null" json:"answer"`
	Explanation       string `gorm:"type:text" json:"explanation"`
}

// ─── Quiz ─────────────────────────────────────────────

type Quiz struct {
	BaseModel
	QuestionID uint   `json:"question_id"`
	UserID     uint   `json:"user_id"`
	UserAnswer string `gorm:"size:20;not null" json:"user_answer"`
	IsCorrect  bool   `json:"is_correct"`
}

// ─── Ask (Q&A) ────────────────────────────────────────

type AskSession struct {
	BaseModel
	UserID  uint   `json:"user_id"`
	Title   string `gorm:"size:200" json:"title"`
}

type AskMessage struct {
	BaseModel
	SessionID   uint    `json:"session_id"`
	Role        string  `gorm:"size:20;not null" json:"role"` // user/assistant
	Content     string  `gorm:"type:text" json:"content"`
	Confidence  float64 `json:"confidence"`
}
