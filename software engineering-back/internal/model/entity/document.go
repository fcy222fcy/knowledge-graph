package entity

// ─── Document ─────────────────────────────────────────

// Document 文档实体，存储上传的文档信息和解析后的内容
type Document struct {
	BaseModel
	UserID        uint   `gorm:"index;comment:用户ID" json:"user_id"`
	Title         string `gorm:"size:200;not null;comment:文档标题" json:"title"`
	Description   string `gorm:"size:500;comment:文档描述" json:"description"`
	Filename      string `gorm:"size:200;not null;comment:原始文件名" json:"filename"`
	FilePath      string `gorm:"size:500;not null;comment:文件存储路径" json:"-"`          // 文件存储路径，不序列化
	FileSize      int64  `gorm:"comment:文件大小字节" json:"file_size"`
	FileType      string `gorm:"size:20;comment:文件类型 pdf/docx/txt 等" json:"file_type"` // 文件类型：pdf/docx/txt 等
	Content       string `gorm:"type:longtext;comment:解析后的文本内容" json:"-"`               // 解析后的文本内容，不序列化
	Status        string `gorm:"size:20;default:pending;comment:处理状态 pending/approved/rejected/completed/failed" json:"status"` // 处理状态：pending/approved/rejected/completed/failed
	ReviewComment string `gorm:"size:500;comment:审核意见" json:"review_comment"` // 审核意见
}
