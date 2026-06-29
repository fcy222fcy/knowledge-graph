package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/repository"
	"software_engineering/pkg/response"
)

// GetAnalyticsOverview 获取系统概览统计
func GetAnalyticsOverview(c *gin.Context) {
	// 获取用户总数
	userCount, err := repository.CountUsers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户统计失败")
		return
	}

	// 获取文档总数
	docCount, err := repository.CountDocuments()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文档统计失败")
		return
	}

	// 获取待审核文档数量
	pendingDocCount, err := repository.CountDocumentsByStatus("pending")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取待审核文档统计失败")
		return
	}

	// 获取知识点总数
	kpCount, err := repository.CountKnowledgePoints()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取知识点统计失败")
		return
	}

	// 获取题目总数
	questionCount, err := repository.CountQuestions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取题目统计失败")
		return
	}

	// 获取问答会话总数
	sessionCount, err := repository.CountAskSessions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取问答统计失败")
		return
	}

	// 获取答题记录总数
	quizCount, err := repository.CountQuizzes()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取答题统计失败")
		return
	}

	// 获取作业总数
	assignmentCount, err := repository.CountAssignments()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取作业统计失败")
		return
	}

	response.Success(c, gin.H{
		"user_count":           userCount,
		"document_count":       docCount,
		"pending_document_count": pendingDocCount,
		"knowledge_count":      kpCount,
		"question_count":       questionCount,
		"session_count":        sessionCount,
		"quiz_count":           quizCount,
		"assignment_count":     assignmentCount,
	})
}

// GetUserStats 获取用户统计信息
func GetUserStats(c *gin.Context) {
	// 统计用户总数
	count, err := repository.CountStudents()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户统计失败")
		return
	}

	response.Success(c, gin.H{
		"student": count,
	})
}
