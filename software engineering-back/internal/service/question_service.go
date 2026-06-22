package service

import (
	"encoding/json"
	"errors"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

// parseOptions 解析题目选项 JSON 字符串
func parseOptions(optionsJSON string) []response.QuestionOption {
	var options []response.QuestionOption
	json.Unmarshal([]byte(optionsJSON), &options)
	return options
}

// CreateQuestion 创建题目
func CreateQuestion(req request.CreateQuestionRequest) (uint, error) {
	optionsJSON, _ := json.Marshal(req.Options)
	q := &entity.Question{
		Title:            req.Title,
		Type:             req.Type,
		Difficulty:       req.Difficulty,
		KnowledgePointID: req.KnowledgePointID,
		Options:          string(optionsJSON),
		Answer:           req.Answer,
		Explanation:      req.Explanation,
	}
	if err := repository.CreateQuestion(q); err != nil {
		return 0, err
	}
	return q.ID, nil
}

// GetQuestion 获取题目详情，可选是否包含答案和解析
func GetQuestion(id uint, includeAnswer bool) (*response.QuestionResponse, error) {
	q, err := repository.FindQuestionByID(id)
	if err != nil {
		return nil, errors.New("题目不存在")
	}
	resp := &response.QuestionResponse{
		ID:               q.ID,
		Title:            q.Title,
		Type:             q.Type,
		Difficulty:       q.Difficulty,
		KnowledgePointID: q.KnowledgePointID,
		Options:          parseOptions(q.Options),
		CreatedAt:        q.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if includeAnswer {
		resp.Answer = q.Answer
		resp.Explanation = q.Explanation
	}
	return resp, nil
}

// UpdateQuestion 更新题目信息，仅更新非空字段
func UpdateQuestion(id uint, req request.UpdateQuestionRequest) error {
	q, err := repository.FindQuestionByID(id)
	if err != nil {
		return errors.New("题目不存在")
	}
	if req.Title != "" {
		q.Title = req.Title
	}
	if req.Type != "" {
		q.Type = req.Type
	}
	if req.Difficulty != "" {
		q.Difficulty = req.Difficulty
	}
	if req.Options != nil {
		optionsJSON, _ := json.Marshal(req.Options)
		q.Options = string(optionsJSON)
	}
	if req.Answer != "" {
		q.Answer = req.Answer
	}
	if req.Explanation != "" {
		q.Explanation = req.Explanation
	}
	return repository.UpdateQuestion(q)
}

// DeleteQuestion 删除题目
func DeleteQuestion(id uint) error {
	_, err := repository.FindQuestionByID(id)
	if err != nil {
		return errors.New("题目不存在")
	}
	return repository.DeleteQuestion(id)
}

// ListQuestions 分页查询题目列表，支持按标题、知识点和难度过滤
func ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]response.QuestionResponse, int64, error) {
	questions, total, err := repository.ListQuestions(page, size, keyword, knowledgePointID, difficulty)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.QuestionResponse, len(questions))
	for i, q := range questions {
		list[i] = response.QuestionResponse{
			ID:               q.ID,
			Title:            q.Title,
			Type:             q.Type,
			Difficulty:       q.Difficulty,
			KnowledgePointID: q.KnowledgePointID,
			Options:          parseOptions(q.Options),
			CreatedAt:        q.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}
