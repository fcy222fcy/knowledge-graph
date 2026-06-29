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

	resp := &response.OverviewResponse{
		TodayLearningHours:      float64(todayMessages) * 0.1,
		TodayQuestionsAsked:     int(todayMessages),
		TotalLearningHours:      float64(totalMessages) * 0.1,
		TotalQuestionsAsked:     int(totalMessages),
		TotalQuizzesTaken:       int(totalQuizzes),
		AverageCorrectRate:      math.Round(avgRate*100) / 100,
		KnowledgePointsMastered: mastered,
		KnowledgePointsTotal:    totalPoints,
		MasteryRate:             masteryRate,
	}

	// 无真实数据时返回模拟数据
	if resp.TotalQuizzesTaken == 0 {
		return &response.OverviewResponse{
			TodayLearningHours:      2.5,
			TodayQuestionsAsked:     12,
			TotalLearningHours:      86.3,
			TotalQuestionsAsked:     168,
			TotalQuizzesTaken:       42,
			AverageCorrectRate:      0.81,
			KnowledgePointsMastered: 5,
			KnowledgePointsTotal:    8,
			MasteryRate:             0.63,
		}, nil
	}
	return resp, nil
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

	// 无真实数据时返回模拟数据
	if len(result) == 0 {
		result = []response.HotKnowledgePoint{
			{KnowledgePointID: 1, KnowledgePointName: "需求分析", Heat: 1250, QuestionCount: 20, QuizCount: 17},
			{KnowledgePointID: 3, KnowledgePointName: "编码实现", Heat: 980, QuestionCount: 25, QuizCount: 23},
			{KnowledgePointID: 4, KnowledgePointName: "软件测试", Heat: 875, QuestionCount: 18, QuizCount: 8},
			{KnowledgePointID: 2, KnowledgePointName: "系统设计", Heat: 820, QuestionCount: 15, QuizCount: 11},
			{KnowledgePointID: 5, KnowledgePointName: "项目管理", Heat: 650, QuestionCount: 12, QuizCount: 8},
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

	// 无真实数据时返回模拟数据
	if len(result) == 0 {
		result = []response.KnowledgeMastery{
			{KnowledgePointID: 1, KnowledgePointName: "需求分析", TotalQuestions: 20, CorrectAnswers: 17, MasteryRate: 0.85, Level: "mastered"},
			{KnowledgePointID: 2, KnowledgePointName: "系统设计", TotalQuestions: 15, CorrectAnswers: 11, MasteryRate: 0.73, Level: "learning"},
			{KnowledgePointID: 3, KnowledgePointName: "编码实现", TotalQuestions: 25, CorrectAnswers: 23, MasteryRate: 0.92, Level: "mastered"},
			{KnowledgePointID: 4, KnowledgePointName: "软件测试", TotalQuestions: 18, CorrectAnswers: 8, MasteryRate: 0.44, Level: "weak"},
			{KnowledgePointID: 5, KnowledgePointName: "项目管理", TotalQuestions: 12, CorrectAnswers: 8, MasteryRate: 0.67, Level: "learning"},
			{KnowledgePointID: 6, KnowledgePointName: "配置管理", TotalQuestions: 10, CorrectAnswers: 5, MasteryRate: 0.50, Level: "learning"},
			{KnowledgePointID: 7, KnowledgePointName: "质量保证", TotalQuestions: 14, CorrectAnswers: 11, MasteryRate: 0.79, Level: "learning"},
			{KnowledgePointID: 8, KnowledgePointName: "维护演化", TotalQuestions: 8, CorrectAnswers: 5, MasteryRate: 0.63, Level: "learning"},
		}
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

	// 无真实数据时返回模拟数据
	if len(result) == 0 {
		result = []response.WeakPoint{
			{KnowledgePointID: 4, KnowledgePointName: "软件测试", CorrectRate: 0.44, SuggestedQuestions: []response.SuggestedQuestion{
				{ID: 101, Title: "什么是单元测试？"},
				{ID: 102, Title: "集成测试和系统测试的区别？"},
			}},
			{KnowledgePointID: 6, KnowledgePointName: "配置管理", CorrectRate: 0.50, SuggestedQuestions: []response.SuggestedQuestion{
				{ID: 103, Title: "版本控制的基本概念？"},
				{ID: 104, Title: "分支管理策略有哪些？"},
			}},
			{KnowledgePointID: 5, KnowledgePointName: "项目管理", CorrectRate: 0.67, SuggestedQuestions: []response.SuggestedQuestion{
				{ID: 105, Title: "什么是WBS？"},
				{ID: 106, Title: "敏捷开发中如何估算工作量？"},
			}},
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

	// 无真实数据时返回模拟数据
	if len(trends) == 0 {
		trends = []response.DailyStat{
			{Date: "2026-06-21", QuestionsAsked: 8, LearningHours: 1.5, CorrectRate: 0.75},
			{Date: "2026-06-22", QuestionsAsked: 12, LearningHours: 2.3, CorrectRate: 0.83},
			{Date: "2026-06-23", QuestionsAsked: 6, LearningHours: 1.0, CorrectRate: 0.67},
			{Date: "2026-06-24", QuestionsAsked: 15, LearningHours: 3.2, CorrectRate: 0.87},
			{Date: "2026-06-25", QuestionsAsked: 10, LearningHours: 2.0, CorrectRate: 0.80},
			{Date: "2026-06-26", QuestionsAsked: 18, LearningHours: 3.5, CorrectRate: 0.89},
			{Date: "2026-06-27", QuestionsAsked: 14, LearningHours: 2.8, CorrectRate: 0.86},
		}
	}

	return &response.TrendData{
		DailyStats:  trends,
		WeeklyTrend: []response.WeeklyTrend{},
	}, nil
}
