package repository

import (
	"time"

	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CountQuizzesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func CountCorrectQuizzesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&count).Error
	return count, err
}

func CountTodayQuizzesByUser(userID uint) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&count).Error
	return count, err
}

func CountTodayMessagesByUser(userID uint) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&model.AskMessage{}).
		Joins("JOIN ask_sessions ON ask_messages.session_id = ask_sessions.id").
		Where("ask_sessions.user_id = ? AND ask_messages.role = 'user' AND DATE(ask_messages.created_at) = ?", userID, today).
		Count(&count).Error
	return count, err
}

func CountTotalMessagesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.AskMessage{}).
		Joins("JOIN ask_sessions ON ask_messages.session_id = ask_sessions.id").
		Where("ask_sessions.user_id = ? AND ask_messages.role = 'user'", userID).
		Count(&count).Error
	return count, err
}

func GetQuizzesByKnowledgePoint(userID uint) (map[uint]int, map[uint]int, error) {
	var results []struct {
		KnowledgePointID uint
		Total            int
		Correct          int
	}
	query := database.DB.Model(&model.Quiz{}).
		Select("questions.knowledge_point_id, COUNT(*) as total, SUM(CASE WHEN quizzes.is_correct = 1 THEN 1 ELSE 0 END) as correct").
		Joins("JOIN questions ON quizzes.question_id = questions.id")
	if userID > 0 {
		query = query.Where("quizzes.user_id = ?", userID)
	}
	err := query.Group("questions.knowledge_point_id").Scan(&results).Error

	totalMap := make(map[uint]int)
	correctMap := make(map[uint]int)
	for _, r := range results {
		totalMap[r.KnowledgePointID] = r.Total
		correctMap[r.KnowledgePointID] = r.Correct
	}
	return totalMap, correctMap, err
}

func GetDailyQuizStats(userID uint, days int) ([]struct {
	Date    string
	Correct int
	Total   int
}, error) {
	var results []struct {
		Date    string
		Correct int
		Total   int
	}
	since := time.Now().AddDate(0, 0, -days)
	err := database.DB.Model(&model.Quiz{}).
		Select("DATE(created_at) as date, COUNT(*) as total, SUM(CASE WHEN is_correct = 1 THEN 1 ELSE 0 END) as correct").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}

func FindQuestionsByKnowledgePoint(kpID uint, limit int) ([]model.Question, error) {
	var questions []model.Question
	err := database.DB.Where("knowledge_point_id = ?", kpID).Limit(limit).Find(&questions).Error
	return questions, err
}
