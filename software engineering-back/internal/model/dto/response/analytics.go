package response

// OverviewResponse 学习概览响应
type OverviewResponse struct {
	TodayLearningHours      float64 `json:"today_learning_hours"`      // 今日学习时长（小时）
	TodayQuestionsAsked     int     `json:"today_questions_asked"`     // 今日提问数量
	TotalLearningHours      float64 `json:"total_learning_hours"`      // 总学习时长（小时）
	TotalQuestionsAsked     int     `json:"total_questions_asked"`     // 总提问数量
	TotalQuizzesTaken       int     `json:"total_quizzes_taken"`       // 总答题次数
	AverageCorrectRate      float64 `json:"average_correct_rate"`      // 平均正确率
	KnowledgePointsMastered int     `json:"knowledge_points_mastered"` // 已掌握知识点数量
	KnowledgePointsTotal    int     `json:"knowledge_points_total"`    // 总知识点数量
	MasteryRate             float64 `json:"mastery_rate"`             // 知识点掌握率
}

// HotKnowledgePoint 热门知识点
type HotKnowledgePoint struct {
	KnowledgePointID   uint   `json:"knowledge_point_id"`   // 知识点ID
	KnowledgePointName string `json:"knowledge_point_name"` // 知识点名称
	Heat               int    `json:"heat"`                 // 热度值
	QuestionCount      int    `json:"question_count"`       // 相关题目数量
	QuizCount          int    `json:"quiz_count"`           // 答题次数
}

// KnowledgeMastery 知识点掌握情况
type KnowledgeMastery struct {
	KnowledgePointID   uint    `json:"knowledge_point_id"`   // 知识点ID
	KnowledgePointName string  `json:"knowledge_point_name"` // 知识点名称
	TotalQuestions      int     `json:"total_questions"`      // 总题目数
	CorrectAnswers      int     `json:"correct_answers"`      // 正确答题数
	MasteryRate         float64 `json:"mastery_rate"`         // 掌握率
	Level               string  `json:"level"`                // 掌握等级（如：初级/中级/高级）
}

// WeakPoint 薄弱知识点
type WeakPoint struct {
	KnowledgePointID   uint                `json:"knowledge_point_id"`   // 知识点ID
	KnowledgePointName string              `json:"knowledge_point_name"` // 知识点名称
	CorrectRate        float64             `json:"correct_rate"`         // 正确率
	SuggestedQuestions []SuggestedQuestion `json:"suggested_questions"`  // 建议练习的题目
}

// SuggestedQuestion 建议练习的题目
type SuggestedQuestion struct {
	ID    uint   `json:"id"`    // 题目ID
	Title string `json:"title"` // 题目标题
}

// TrendData 学习趋势数据
type TrendData struct {
	DailyStats  []DailyStat   `json:"daily_stats"`  // 每日统计
	WeeklyTrend []WeeklyTrend `json:"weekly_trend"` // 每周趋势
}

// DailyStat 每日学习统计
type DailyStat struct {
	Date           string  `json:"date"`            // 日期
	QuestionsAsked int     `json:"questions_asked"` // 提问数量
	LearningHours  float64 `json:"learning_hours"`  // 学习时长（小时）
	CorrectRate    float64 `json:"correct_rate"`    // 正确率
}

// WeeklyTrend 每周学习趋势
type WeeklyTrend struct {
	Week                string  `json:"week"`                 // 周次
	AvgCorrectRate      float64 `json:"avg_correct_rate"`     // 平均正确率
	TotalLearningHours  float64 `json:"total_learning_hours"` // 总学习时长
	TotalQuestionsAsked int     `json:"total_questions_asked"` // 总提问数量
}
