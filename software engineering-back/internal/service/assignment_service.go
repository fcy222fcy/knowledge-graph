package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

// ─── 教师端 ────────────────────────────────────────────

// CreateAssignment 创建作业（含题目）
func CreateAssignment(teacherID uint, req request.CreateAssignmentRequest) (*response.AssignmentResponse, error) {
	// 计算总分
	totalScore := 0
	for _, q := range req.Questions {
		if q.Score > 0 {
			totalScore += q.Score
		} else {
			totalScore += 10 // 默认10分
		}
	}

	assignment := &entity.Assignment{
		Title:       req.Title,
		Description: req.Description,
		Chapter:     req.Chapter,
		Deadline:    req.Deadline,
		Status:      "draft",
		TotalScore:  totalScore,
		TeacherID:   teacherID,
	}
	if err := repository.CreateAssignment(assignment); err != nil {
		return nil, fmt.Errorf("创建作业失败: %w", err)
	}

	// 创建题目
	questions := make([]entity.AssignmentQuestion, 0, len(req.Questions))
	for i, q := range req.Questions {
		optionsJSON, _ := json.Marshal(q.Options)
		score := q.Score
		if score <= 0 {
			score = 10
		}
		questions = append(questions, entity.AssignmentQuestion{
			AssignmentID: assignment.ID,
			Title:        q.Title,
			Type:         q.Type,
			Options:      string(optionsJSON),
			Answer:       q.Answer,
			Explanation:  q.Explanation,
			Score:        score,
			SortOrder:    q.SortOrder,
		})
		if questions[i].SortOrder == 0 {
			questions[i].SortOrder = i + 1
		}
	}
	if err := repository.CreateAssignmentQuestions(questions); err != nil {
		return nil, fmt.Errorf("创建题目失败: %w", err)
	}

	return &response.AssignmentResponse{
		ID:          assignment.ID,
		Title:       assignment.Title,
		Chapter:     assignment.Chapter,
		Deadline:    assignment.Deadline,
		Status:      assignment.Status,
		TotalScore:  assignment.TotalScore,
		QuestionNum: len(questions),
		TeacherID:   assignment.TeacherID,
		CreatedAt:   assignment.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateAssignment 更新作业
func UpdateAssignment(assignmentID uint, req request.UpdateAssignmentRequest) error {
	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return fmt.Errorf("作业不存在")
	}
	if assignment.Status == "closed" {
		return fmt.Errorf("已截止的作业不能编辑")
	}

	if req.Title != "" {
		assignment.Title = req.Title
	}
	if req.Description != "" {
		assignment.Description = req.Description
	}
	if req.Chapter != "" {
		assignment.Chapter = req.Chapter
	}
	if req.Deadline != "" {
		assignment.Deadline = req.Deadline
	}

	// 如果提供了题目列表，则替换
	if req.Questions != nil {
		repository.DeleteAssignmentQuestions(assignmentID)

		totalScore := 0
		questions := make([]entity.AssignmentQuestion, 0, len(req.Questions))
		for i, q := range req.Questions {
			optionsJSON, _ := json.Marshal(q.Options)
			score := q.Score
			if score <= 0 {
				score = 10
			}
			totalScore += score
			questions = append(questions, entity.AssignmentQuestion{
				AssignmentID: assignmentID,
				Title:        q.Title,
				Type:         q.Type,
				Options:      string(optionsJSON),
				Answer:       q.Answer,
				Explanation:  q.Explanation,
				Score:        score,
				SortOrder:    q.SortOrder,
			})
			if questions[i].SortOrder == 0 {
				questions[i].SortOrder = i + 1
			}
		}
		repository.CreateAssignmentQuestions(questions)
		assignment.TotalScore = totalScore
	}

	return repository.UpdateAssignment(assignment)
}

// PublishAssignment 发布作业
func PublishAssignment(assignmentID uint) error {
	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return fmt.Errorf("作业不存在")
	}
	assignment.Status = "published"
	return repository.UpdateAssignment(assignment)
}

// CloseAssignment 关闭作业
func CloseAssignment(assignmentID uint) error {
	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return fmt.Errorf("作业不存在")
	}
	assignment.Status = "closed"
	return repository.UpdateAssignment(assignment)
}

// DeleteAssignment 删除作业
func DeleteAssignment(assignmentID uint) error {
	repository.DeleteAssignmentQuestions(assignmentID)
	return repository.DeleteAssignment(assignmentID)
}

// GetAssignmentDetail 教师查看作业详情（含答案）
func GetAssignmentDetail(assignmentID uint) (*response.AssignmentDetailResponse, error) {
	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return nil, fmt.Errorf("作业不存在")
	}

	questions, err := repository.ListAssignmentQuestions(assignmentID)
	if err != nil {
		return nil, err
	}

	return buildAssignmentDetail(assignment, questions, true), nil
}

// ListTeacherAssignments 教师查看作业列表
func ListTeacherAssignments(teacherID uint, page, size int) ([]response.AssignmentResponse, int64, error) {
	list, total, err := repository.ListAssignmentsByTeacher(teacherID, page, size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.AssignmentResponse, len(list))
	for i, a := range list {
		result[i] = response.AssignmentResponse{
			ID:          a.ID,
			Title:       a.Title,
			Chapter:     a.Chapter,
			Deadline:    a.Deadline,
			Status:      a.Status,
			TotalScore:  a.TotalScore,
			QuestionNum: repository.CountAssignmentQuestions(a.ID),
			SubmitCount: repository.CountAssignmentSubmissions(a.ID),
			TeacherID:   a.TeacherID,
			CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return result, total, nil
}

// ListAssignmentSubmissions 查看作业提交列表
func ListAssignmentSubmissions(assignmentID uint, page, size int) ([]response.AssignmentSubmissionResponse, int64, error) {
	list, total, err := repository.ListAssignmentSubmissions(assignmentID, page, size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.AssignmentSubmissionResponse, len(list))
	for i, s := range list {
		// 查询用户名
		user, _ := repository.FindUserByID(s.UserID)
		username := ""
		if user != nil {
			username = user.Username
		}
		result[i] = response.AssignmentSubmissionResponse{
			ID:           s.ID,
			AssignmentID: s.AssignmentID,
			UserID:       s.UserID,
			Username:     username,
			Score:        s.Score,
			TotalScore:   s.TotalScore,
			Status:       s.Status,
			SubmittedAt:  s.SubmittedAt,
		}
	}
	return result, total, nil
}

// ─── 学生端 ────────────────────────────────────────────

// ListStudentAssignments 学生查看作业列表
func ListStudentAssignments(userID uint, page, size int) ([]response.AssignmentResponse, int64, error) {
	list, total, err := repository.ListPublishedAssignments(page, size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.AssignmentResponse, len(list))
	for i, a := range list {
		submitCount := repository.CountAssignmentSubmissions(a.ID)
		totalCount := submitCount

		resp := response.AssignmentResponse{
			ID:          a.ID,
			Title:       a.Title,
			Chapter:     a.Chapter,
			Deadline:    a.Deadline,
			Status:      a.Status,
			TotalScore:  a.TotalScore,
			QuestionNum: repository.CountAssignmentQuestions(a.ID),
			SubmitCount: submitCount,
			TotalCount:  totalCount,
			TeacherID:   a.TeacherID,
			CreatedAt:   a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// 查询学生是否已提交
		submission, err := repository.FindAssignmentSubmission(a.ID, userID)
		if err == nil && submission.ID > 0 {
			resp.IsSubmitted = true
			resp.Score = &submission.Score
		}

		result[i] = resp
	}
	return result, total, nil
}

// GetStudentAssignmentDetail 学生查看作业详情（不含答案）
func GetStudentAssignmentDetail(assignmentID uint) (*response.AssignmentDetailResponse, error) {
	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return nil, fmt.Errorf("作业不存在")
	}
	if assignment.Status != "published" {
		return nil, fmt.Errorf("作业未发布")
	}

	questions, err := repository.ListAssignmentQuestions(assignmentID)
	if err != nil {
		return nil, err
	}

	return buildAssignmentDetail(assignment, questions, false), nil
}

// SubmitAssignment 学生提交作业
func SubmitAssignment(userID, assignmentID uint, req request.SubmitAssignmentRequest) (*response.AssignmentSubmitResult, error) {
	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return nil, fmt.Errorf("作业不存在")
	}
	if assignment.Status != "published" {
		return nil, fmt.Errorf("作业未发布或已截止")
	}

	// 检查是否已提交
	existing, err := repository.FindAssignmentSubmission(assignmentID, userID)
	if err == nil && existing.ID > 0 {
		return nil, fmt.Errorf("已提交过该作业")
	}

	// 获取题目用于自动批改
	questions, err := repository.ListAssignmentQuestions(assignmentID)
	if err != nil {
		return nil, err
	}
	questionMap := make(map[uint]entity.AssignmentQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	// 自动批改
	score := 0
	answersData := make([]map[string]interface{}, 0, len(req.Answers))
	for _, ans := range req.Answers {
		q, ok := questionMap[ans.QuestionID]
		isCorrect := false
		if ok && q.Answer == ans.Answer {
			isCorrect = true
			score += q.Score
		}
		answersData = append(answersData, map[string]interface{}{
			"question_id": ans.QuestionID,
			"answer":      ans.Answer,
			"is_correct":  isCorrect,
		})
	}

	answersJSON, _ := json.Marshal(answersData)
	now := time.Now().Format("2006-01-02T15:04:05Z")

	submission := &entity.AssignmentSubmission{
		AssignmentID: assignmentID,
		UserID:       userID,
		Answers:      string(answersJSON),
		Score:        score,
		TotalScore:   assignment.TotalScore,
		Status:       "submitted",
		SubmittedAt:  now,
	}
	if err := repository.CreateAssignmentSubmission(submission); err != nil {
		return nil, fmt.Errorf("提交失败: %w", err)
	}

	log.Printf("学生 %d 提交作业 %d，得分 %d/%d", userID, assignmentID, score, assignment.TotalScore)

	return &response.AssignmentSubmitResult{
		SubmissionID: submission.ID,
		Score:        score,
		TotalScore:   assignment.TotalScore,
		Status:       submission.Status,
	}, nil
}

// GetStudentAssignmentResult 学生查看作业结果（含题目和答题情况）
func GetStudentAssignmentResult(userID, assignmentID uint) (*response.AssignmentResultResponse, error) {
	submission, err := repository.FindAssignmentSubmission(assignmentID, userID)
	if err != nil {
		return nil, fmt.Errorf("未提交该作业")
	}

	assignment, err := repository.FindAssignmentByID(assignmentID)
	if err != nil {
		return nil, fmt.Errorf("作业不存在")
	}

	questions, err := repository.ListAssignmentQuestions(assignmentID)
	if err != nil {
		return nil, err
	}

	// 解析学生的答题记录
	var answersData []map[string]interface{}
	json.Unmarshal([]byte(submission.Answers), &answersData)
	studentAnswers := make(map[uint]map[string]interface{})
	for _, a := range answersData {
		if qid, ok := a["question_id"].(float64); ok {
			studentAnswers[uint(qid)] = a
		}
	}

	// 组装题目结果
	qResults := make([]response.AssignmentQuestionResult, len(questions))
	for i, q := range questions {
		var options []response.QuestionOption
		if q.Options != "" {
			json.Unmarshal([]byte(q.Options), &options)
		}
		item := response.AssignmentQuestionResult{
			ID:          q.ID,
			Title:       q.Title,
			Type:        q.Type,
			Options:     options,
			Score:       q.Score,
			SortOrder:   q.SortOrder,
			Answer:      q.Answer,
			Explanation: q.Explanation,
		}
		if sa, ok := studentAnswers[q.ID]; ok {
			if ans, _ := sa["answer"].(string); ans != "" {
				item.MyAnswer = ans
			}
			if correct, _ := sa["is_correct"].(bool); correct {
				item.IsCorrect = true
			}
		}
		qResults[i] = item
	}

	return &response.AssignmentResultResponse{
		ID:           submission.ID,
		AssignmentID: submission.AssignmentID,
		Title:        assignment.Title,
		Score:        submission.Score,
		TotalScore:   submission.TotalScore,
		Status:       submission.Status,
		SubmittedAt:  submission.SubmittedAt,
		Questions:    qResults,
	}, nil
}

// ─── 辅助函数 ──────────────────────────────────────────

func buildAssignmentDetail(assignment *entity.Assignment, questions []entity.AssignmentQuestion, includeAnswer bool) *response.AssignmentDetailResponse {
	qList := make([]response.AssignmentQuestionResponse, len(questions))
	for i, q := range questions {
		var options []response.QuestionOption
		if q.Options != "" {
			json.Unmarshal([]byte(q.Options), &options)
		}
		item := response.AssignmentQuestionResponse{
			ID:        q.ID,
			Title:     q.Title,
			Type:      q.Type,
			Options:   options,
			Score:     q.Score,
			SortOrder: q.SortOrder,
		}
		if includeAnswer {
			item.Answer = q.Answer
			item.Explanation = q.Explanation
		}
		qList[i] = item
	}

	return &response.AssignmentDetailResponse{
		ID:          assignment.ID,
		Title:       assignment.Title,
		Description: assignment.Description,
		Chapter:     assignment.Chapter,
		Deadline:    assignment.Deadline,
		TotalScore:  assignment.TotalScore,
		Status:      assignment.Status,
		Questions:   qList,
	}
}
