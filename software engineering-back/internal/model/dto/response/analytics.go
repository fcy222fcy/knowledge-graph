package response

type OverviewResponse struct {
	TodayLearningHours      float64 `json:"today_learning_hours"`
	TodayQuestionsAsked     int     `json:"today_questions_asked"`
	TotalLearningHours      float64 `json:"total_learning_hours"`
	TotalQuestionsAsked     int     `json:"total_questions_asked"`
	TotalQuizzesTaken       int     `json:"total_quizzes_taken"`
	AverageCorrectRate      float64 `json:"average_correct_rate"`
	KnowledgePointsMastered int     `json:"knowledge_points_mastered"`
	KnowledgePointsTotal    int     `json:"knowledge_points_total"`
	MasteryRate             float64 `json:"mastery_rate"`
}

type HotKnowledgePoint struct {
	KnowledgePointID   uint   `json:"knowledge_point_id"`
	KnowledgePointName string `json:"knowledge_point_name"`
	Heat               int    `json:"heat"`
	QuestionCount      int    `json:"question_count"`
	QuizCount          int    `json:"quiz_count"`
}

type KnowledgeMastery struct {
	KnowledgePointID   uint    `json:"knowledge_point_id"`
	KnowledgePointName string  `json:"knowledge_point_name"`
	TotalQuestions      int     `json:"total_questions"`
	CorrectAnswers      int     `json:"correct_answers"`
	MasteryRate         float64 `json:"mastery_rate"`
	Level               string  `json:"level"`
}

type WeakPoint struct {
	KnowledgePointID   uint                `json:"knowledge_point_id"`
	KnowledgePointName string              `json:"knowledge_point_name"`
	CorrectRate        float64             `json:"correct_rate"`
	SuggestedQuestions []SuggestedQuestion `json:"suggested_questions"`
}

type SuggestedQuestion struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type TrendData struct {
	DailyStats  []DailyStat   `json:"daily_stats"`
	WeeklyTrend []WeeklyTrend `json:"weekly_trend"`
}

type DailyStat struct {
	Date           string  `json:"date"`
	QuestionsAsked int     `json:"questions_asked"`
	LearningHours  float64 `json:"learning_hours"`
	CorrectRate    float64 `json:"correct_rate"`
}

type WeeklyTrend struct {
	Week                string  `json:"week"`
	AvgCorrectRate      float64 `json:"avg_correct_rate"`
	TotalLearningHours  float64 `json:"total_learning_hours"`
	TotalQuestionsAsked int     `json:"total_questions_asked"`
}
