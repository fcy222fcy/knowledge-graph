package service

import (
	"fmt"
	"strings"

	"software_engineering/internal/model/dto"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

func CreateSession(userID uint, req dto.CreateSessionRequest) (*dto.AskSessionResponse, error) {
	session := &entity.AskSession{
		UserID: userID,
		Title:  req.Title,
	}
	if err := repository.CreateAskSession(session); err != nil {
		return nil, err
	}
	return &dto.AskSessionResponse{
		ConversationID: session.ID,
		Title:          session.Title,
		UpdatedAt:      session.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListSessions(userID uint, page, size int) ([]dto.AskSessionResponse, int64, error) {
	sessions, total, err := repository.ListAskSessionsByUser(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.AskSessionResponse, len(sessions))
	for i, s := range sessions {
		item := dto.AskSessionResponse{
			ConversationID: s.ID,
			Title:          s.Title,
			MessageCount:   repository.CountMessagesBySession(s.ID),
			UpdatedAt:      s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
		lastMsg, err := repository.GetLastMessageBySession(s.ID)
		if err == nil {
			item.LastQuestion = lastMsg.Content
		}
		list[i] = item
	}
	return list, total, nil
}

func ListSessionMessages(sessionID uint, page, size int) ([]dto.AskMessageResponse, int64, error) {
	messages, total, err := repository.ListAskMessages(sessionID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.AskMessageResponse, len(messages))
	for i, m := range messages {
		list[i] = dto.AskMessageResponse{
			MessageID: m.ID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

func Ask(userID uint, req dto.AskRequest) (*dto.AskResponse, error) {
	// Auto-create or reuse session
	sessionID := req.ConversationID
	if sessionID == 0 {
		session := &entity.AskSession{
			UserID: userID,
			Title:  req.Question,
		}
		repository.CreateAskSession(session)
		sessionID = session.ID
	}

	// Save user message
	userMsg := &entity.AskMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Question,
	}
	repository.CreateAskMessage(userMsg)

	// Simplified keyword-matching answer
	answer := fmt.Sprintf("关于「%s」的回答：这是系统生成的示例回答。在生产环境中，这里会接入语义检索和 AI 生成。", req.Question)

	// Find related knowledge points
	points, _ := repository.GetAllKnowledgePoints()
	var related []dto.KPRef
	for _, p := range points {
		if strings.Contains(req.Question, p.Name) || strings.Contains(p.Name, req.Question) {
			related = append(related, dto.KPRef{ID: p.ID, Name: p.Name, Description: p.Description})
		}
	}

	// Save assistant message
	assistantMsg := &entity.AskMessage{
		SessionID:  sessionID,
		Role:       "assistant",
		Content:    answer,
		Confidence: 0.75,
	}
	repository.CreateAskMessage(assistantMsg)

	return &dto.AskResponse{
		ConversationID:         sessionID,
		QuestionID:             userMsg.ID,
		Answer:                 answer,
		Confidence:             0.75,
		Sources:                []dto.AskSource{},
		RelatedKnowledgePoints: related,
		CreatedAt:              assistantMsg.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListAskHistory(userID uint, page, size int, conversationID uint) ([]dto.AskResponse, int64, error) {
	// Simplified: return session list as history
	sessions, total, err := repository.ListAskSessionsByUser(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.AskResponse, len(sessions))
	for i, s := range sessions {
		lastMsg, _ := repository.GetLastMessageBySession(s.ID)
		question := ""
		if lastMsg.ID > 0 {
			question = lastMsg.Content
		}
		list[i] = dto.AskResponse{
			ConversationID: s.ID,
			Answer:         question,
			CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}
