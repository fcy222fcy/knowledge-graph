package entity

// ─── Question ─────────────────────────────────────────

// Question 题目实体，存储试题信息
type Question struct {
	BaseModel
	Title            string `gorm:"size:500;not null;comment:题目标题" json:"title"`           // 题目标题
	Type             string `gorm:"size:20;not null;comment:题目类型 single/multiple" json:"type"`             // 题目类型：single/multiple
	Difficulty       string `gorm:"size:20;not null;comment:难度 easy/medium/hard" json:"difficulty"`       // 难度：easy/medium/hard
	KnowledgePointID uint   `gorm:"comment:关联的知识点ID" json:"knowledge_point_id"`                       // 关联的知识点 ID
	Options          string `gorm:"type:text;comment:选项 JSON数组" json:"-"`                       // 选项，JSON 数组存储为文本
	Answer           string `gorm:"size:20;not null;comment:正确答案" json:"answer"`           // 正确答案
	Explanation      string `gorm:"type:text;comment:题目解析" json:"explanation"`             // 题目解析
}
