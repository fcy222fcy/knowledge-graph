package repository

import (
	"time"

	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CountQuizzesByUser 统计用户答题总数
func CountQuizzesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountCorrectQuizzesByUser 统计用户答对的题目数量
func CountCorrectQuizzesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&count).Error
	return count, err
}

// CountTodayQuizzesByUser 统计用户今日答题数量
func CountTodayQuizzesByUser(userID uint) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&entity.Quiz{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&count).Error
	return count, err
}

// CountTodayMessagesByUser 统计用户今日提问数量
func CountTodayMessagesByUser(userID uint) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&entity.AskMessage{}).
		Joins("JOIN ask_sessions ON ask_messages.session_id = ask_sessions.id").
		Where("ask_sessions.user_id = ? AND ask_messages.role = 'user' AND DATE(ask_messages.created_at) = ?", userID, today).
		Count(&count).Error
	return count, err
}

// CountTotalMessagesByUser 统计用户历史提问总数
func CountTotalMessagesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.AskMessage{}).
		Joins("JOIN ask_sessions ON ask_messages.session_id = ask_sessions.id").
		Where("ask_sessions.user_id = ? AND ask_messages.role = 'user'", userID).
		Count(&count).Error
	return count, err
}

// GetQuizzesByKnowledgePoint 按知识点统计用户答题情况，返回各知识点的总数和正确数
func GetQuizzesByKnowledgePoint(userID uint) (map[uint]int, map[uint]int, error) {
	var results []struct {
		KnowledgePointID uint
		Total            int
		Correct          int
	}
	query := database.DB.Model(&entity.Quiz{}).
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

// GetDailyQuizStats 获取用户近 N 天的每日答题统计
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
	err := database.DB.Model(&entity.Quiz{}).
		Select("DATE(created_at) as date, COUNT(*) as total, SUM(CASE WHEN is_correct = 1 THEN 1 ELSE 0 END) as correct").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}

// FindQuestionsByKnowledgePoint 根据知识点 ID 获取题目列表
func FindQuestionsByKnowledgePoint(kpID uint, limit int) ([]entity.Question, error) {
	var questions []entity.Question
	err := database.DB.Where("knowledge_point_id = ?", kpID).Limit(limit).Find(&questions).Error
	return questions, err
}
