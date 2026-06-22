package service

import (
	"math"

	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/repository"
)

// GetOverview 获取用户学习概览数据，包括答题统计、知识点掌握率等
func GetOverview(userID uint) (*response.OverviewResponse, error) {
	totalQuizzes, _ := repository.CountQuizzesByUser(userID)
	correctQuizzes, _ := repository.CountCorrectQuizzesByUser(userID)
	_, _ = repository.CountTodayQuizzesByUser(userID)
	todayMessages, _ := repository.CountTodayMessagesByUser(userID)
	totalMessages, _ := repository.CountTotalMessagesByUser(userID)

	points, _ := repository.GetAllKnowledgePointsForGraph()
	totalPoints := len(points)

	// Knowledge points with >80% correct rate are considered mastered
	totalMap, correctMap, _ := repository.GetQuizzesByKnowledgePoint(userID)
	mastered := 0
	for kpID, total := range totalMap {
		correct := correctMap[kpID]
		if total > 0 && float64(correct)/float64(total) >= 0.8 {
			mastered++
		}
	}

	var avgRate float64
	if totalQuizzes > 0 {
		avgRate = float64(correctQuizzes) / float64(totalQuizzes)
	}

	var masteryRate float64
	if totalPoints > 0 {
		masteryRate = math.Round(float64(mastered)/float64(totalPoints)*100) / 100
	}

	return &response.OverviewResponse{
		TodayLearningHours:      float64(todayMessages) * 0.1,
		TodayQuestionsAsked:     int(todayMessages),
		TotalLearningHours:      float64(totalMessages) * 0.1,
		TotalQuestionsAsked:     int(totalMessages),
		TotalQuizzesTaken:       int(totalQuizzes),
		AverageCorrectRate:      math.Round(avgRate*100) / 100,
		KnowledgePointsMastered: mastered,
		KnowledgePointsTotal:    totalPoints,
		MasteryRate:             masteryRate,
	}, nil
}

// GetHotKnowledgePoints 获取热门知识点列表，按答题热度排序
func GetHotKnowledgePoints(limit int) ([]response.HotKnowledgePoint, error) {
	points, _ := repository.GetAllKnowledgePointsForGraph()
	totalMap, correctMap, _ := repository.GetQuizzesByKnowledgePoint(0)

	var result []response.HotKnowledgePoint
	for _, p := range points {
		heat := totalMap[p.ID] * 10
		if heat > 0 {
			result = append(result, response.HotKnowledgePoint{
				KnowledgePointID:   p.ID,
				KnowledgePointName: p.Name,
				Heat:               heat,
				QuestionCount:      totalMap[p.ID],
				QuizCount:          correctMap[p.ID],
			})
		}
	}

	// Sort by heat descending (bubble sort)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Heat > result[i].Heat {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

// GetKnowledgeMastery 获取用户各知识点的掌握程度
func GetKnowledgeMastery(userID uint) ([]response.KnowledgeMastery, error) {
	points, _ := repository.GetAllKnowledgePointsForGraph()
	totalMap, correctMap, _ := repository.GetQuizzesByKnowledgePoint(userID)

	var result []response.KnowledgeMastery
	for _, p := range points {
		total := totalMap[p.ID]
		correct := correctMap[p.ID]
		if total == 0 {
			continue
		}
		rate := float64(correct) / float64(total)
		level := "weak"
		if rate >= 0.8 {
			level = "mastered"
		} else if rate >= 0.5 {
			level = "learning"
		}
		result = append(result, response.KnowledgeMastery{
			KnowledgePointID:   p.ID,
			KnowledgePointName: p.Name,
			TotalQuestions:      total,
			CorrectAnswers:      correct,
			MasteryRate:         math.Round(rate*100) / 100,
			Level:               level,
		})
	}
	return result, nil
}

// GetWeakPoints 获取用户的薄弱知识点，附带推荐练习题目
func GetWeakPoints(userID uint, limit int) ([]response.WeakPoint, error) {
	masteries, _ := GetKnowledgeMastery(userID)

	var result []response.WeakPoint
	for _, m := range masteries {
		if m.Level == "weak" || m.Level == "learning" {
			// Fetch suggested questions for this knowledge point
			suggested := make([]response.SuggestedQuestion, 0)
			questions, _ := repository.FindQuestionsByKnowledgePoint(m.KnowledgePointID, 3)
			for _, q := range questions {
				suggested = append(suggested, response.SuggestedQuestion{
					ID:    q.ID,
					Title: q.Title,
				})
			}
			result = append(result, response.WeakPoint{
				KnowledgePointID:   m.KnowledgePointID,
				KnowledgePointName: m.KnowledgePointName,
				CorrectRate:        m.MasteryRate,
				SuggestedQuestions: suggested,
			})
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

// GetTrends 获取用户近 N 天的学习趋势数据
func GetTrends(userID uint, days int) (*response.TrendData, error) {
	dailyStats, _ := repository.GetDailyQuizStats(userID, days)

	var trends []response.DailyStat
	for _, d := range dailyStats {
		rate := 0.0
		if d.Total > 0 {
			rate = float64(d.Correct) / float64(d.Total)
		}
		trends = append(trends, response.DailyStat{
			Date:           d.Date,
			QuestionsAsked: d.Total,
			LearningHours:  float64(d.Total) * 0.1,
			CorrectRate:    math.Round(rate*100) / 100,
		})
	}

	return &response.TrendData{
		DailyStats:  trends,
		WeeklyTrend: []response.WeeklyTrend{},
	}, nil
}
