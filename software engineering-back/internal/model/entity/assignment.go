package entity

// ─── Assignment ────────────────────────────────────────

// Assignment 作业实体
type Assignment struct {
	BaseModel
	Title       string `gorm:"size:200;not null;comment:作业名称" json:"title"`
	Description string `gorm:"type:text;comment:作业说明" json:"description"`
	Chapter     string `gorm:"size:100;comment:所属章节" json:"chapter"`
	Deadline    string `gorm:"size:30;comment:截止时间" json:"deadline"`
	Status      string `gorm:"size:20;not null;default:draft;comment:状态 draft/published/closed" json:"status"`
	TotalScore  int    `gorm:"default:0;comment:总分" json:"total_score"`
	TeacherID   uint   `gorm:"comment:创建教师ID" json:"teacher_id"`
}

// AssignmentQuestion 作业题目实体
type AssignmentQuestion struct {
	BaseModel
	AssignmentID uint   `gorm:"not null;comment:作业ID" json:"assignment_id"`
	Title        string `gorm:"size:500;not null;comment:题目内容" json:"title"`
	Type         string `gorm:"size:20;not null;comment:题目类型 single/multiple/judge" json:"type"`
	Options      string `gorm:"type:text;comment:选项 JSON数组" json:"-"`
	Answer       string `gorm:"size:20;not null;comment:正确答案" json:"answer"`
	Explanation  string `gorm:"type:text;comment:题目解析" json:"explanation"`
	Score        int    `gorm:"default:10;comment:分值" json:"score"`
	SortOrder    int    `gorm:"default:0;comment:排序" json:"sort_order"`
}

// AssignmentSubmission 作业提交实体
type AssignmentSubmission struct {
	BaseModel
	AssignmentID uint   `gorm:"not null;comment:作业ID" json:"assignment_id"`
	UserID       uint   `gorm:"not null;comment:用户ID" json:"user_id"`
	Answers      string `gorm:"type:text;comment:学生答案 JSON" json:"-"`
	Score        int    `gorm:"default:0;comment:得分" json:"score"`
	TotalScore   int    `gorm:"default:0;comment:满分" json:"total_score"`
	Status       string `gorm:"size:20;default:submitted;comment:状态 submitted/graded" json:"status"`
	SubmittedAt  string `gorm:"size:30;comment:提交时间" json:"submitted_at"`
}
