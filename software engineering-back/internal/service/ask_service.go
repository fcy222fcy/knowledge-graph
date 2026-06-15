package service

import (
	"fmt"
	"log"
	"strings"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

func CreateSession(userID uint, req request.CreateSessionRequest) (*response.AskSessionResponse, error) {
	session := &entity.AskSession{
		UserID: userID,
		Title:  req.Title,
	}
	if err := repository.CreateAskSession(session); err != nil {
		return nil, err
	}
	return &response.AskSessionResponse{
		ConversationID: session.ID,
		Title:          session.Title,
		UpdatedAt:      session.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListSessions(userID uint, page, size int) ([]response.AskSessionResponse, int64, error) {
	sessions, total, err := repository.ListAskSessionsByUser(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.AskSessionResponse, len(sessions))
	for i, s := range sessions {
		item := response.AskSessionResponse{
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

func ListSessionMessages(sessionID uint, page, size int) ([]response.AskMessageResponse, int64, error) {
	messages, total, err := repository.ListAskMessages(sessionID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.AskMessageResponse, len(messages))
	for i, m := range messages {
		list[i] = response.AskMessageResponse{
			MessageID: m.ID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

func Ask(userID uint, req request.AskRequest) (*response.AskResponse, error) {
	// 自动创建或复用会话
	sessionID := req.ConversationID
	if sessionID == 0 {
		session := &entity.AskSession{
			UserID: userID,
			Title:  req.Question,
		}
		repository.CreateAskSession(session)
		sessionID = session.ID
	}

	// 获取历史消息用于上下文
	historyMsgs, _ := repository.ListRecentMessages(sessionID, 10)
	history := make([]ChatMessage, len(historyMsgs))
	for i, m := range historyMsgs {
		history[i] = ChatMessage{Role: m.Role, Content: m.Content}
	}

	// 保存当前用户消息
	userMsg := &entity.AskMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Question,
	}
	repository.CreateAskMessage(userMsg)

	// 尝试调用 AI 服务（带历史上下文）
	var answer string
	var confidence float64
	var sources []response.AskSource
	var related []response.KPRef

	if aiClient.IsAvailable() {
		answerResp, err := aiClient.SearchAndAnswerWithHistory(req.Question, history, 3)
		if err == nil && answerResp.Answer != "" {
			answer = answerResp.Answer
			confidence = answerResp.Confidence
			for _, s := range answerResp.Sources {
				sources = append(sources, response.AskSource{
					DocumentID:    uint(s.DocumentID),
					DocumentTitle: s.DocumentTitle,
					Content:       s.Content,
				})
			}
		} else {
			log.Printf("warning: AI search_and_answer failed, degrading to simple search: %v", err)
			searchResp, err := aiClient.Search(req.Question, 3)
			if err == nil && len(searchResp.Results) > 0 {
				answer = fmt.Sprintf("关于「%s」的回答：\n\n", req.Question)
				for i, r := range searchResp.Results {
					answer += fmt.Sprintf("%d. %s\n\n", i+1, r.ChunkText[:min(len(r.ChunkText), 150)])
					sources = append(sources, response.AskSource{
						DocumentID:    uint(r.DocumentID),
						DocumentTitle: fmt.Sprintf("文档 #%d", r.DocumentID),
						Content:       r.ChunkText[:min(len(r.ChunkText), 200)],
					})
				}
				answer += "以上内容来自知识库语义检索，仅供参考。"
				confidence = 0.75
			}
		}
	} else {
		log.Println("info: AI service not available, using local keyword search")
	}

	// 降级到本地关键词检索
	if answer == "" {
		docs, _ := repository.GetAllDocumentsContent()
		questionLower := strings.ToLower(req.Question)
		for _, doc := range docs {
			contentLower := strings.ToLower(doc.Content)
			if strings.Contains(contentLower, questionLower) {
				idx := strings.Index(contentLower, questionLower)
				start := max(0, idx-100)
				end := min(len(doc.Content), idx+200)
				snippet := doc.Content[start:end]
				if start > 0 {
					snippet = "..." + snippet
				}
				if end < len(doc.Content) {
					snippet = snippet + "..."
				}
				answer = fmt.Sprintf("关于「%s」的回答：\n\n根据文档《%s》中的内容：\n\n%s\n\n以上内容来自知识库文档检索。", req.Question, doc.Title, snippet)
				sources = append(sources, response.AskSource{
					DocumentID:    doc.ID,
					DocumentTitle: doc.Title,
					Content:       snippet[:min(len(snippet), 200)],
				})
				confidence = 0.75
				break
			}
		}
	}

	// 最终降级：知识点匹配
	if answer == "" {
		points, _ := repository.GetAllKnowledgePoints()
		for _, p := range points {
			if strings.Contains(req.Question, p.Name) || strings.Contains(p.Name, req.Question) {
				related = append(related, response.KPRef{ID: p.ID, Name: p.Name, Description: p.Description})
			}
		}
		if len(related) > 0 {
			answer = fmt.Sprintf("关于「%s」的回答：\n\n", req.Question)
			for i, kp := range related {
				answer += fmt.Sprintf("%d. %s: %s\n\n", i+1, kp.Name, kp.Description)
			}
			answer += "以上内容来自知识点库，仅供参考。"
			confidence = 0.7
		} else {
			answer = fmt.Sprintf("抱歉，暂时无法找到关于「%s」的准确回答。您可以尝试：\n1. 上传更多相关文档\n2. 构建知识图谱\n3. 联系管理员获取帮助", req.Question)
			confidence = 0.3
		}
	}

	// 保存助手消息
	assistantMsg := &entity.AskMessage{
		SessionID:  sessionID,
		Role:       "assistant",
		Content:    answer,
		Confidence: confidence,
	}
	repository.CreateAskMessage(assistantMsg)

	return &response.AskResponse{
		ConversationID:         sessionID,
		QuestionID:             userMsg.ID,
		Answer:                 answer,
		Confidence:             confidence,
		Sources:                sources,
		RelatedKnowledgePoints: related,
		CreatedAt:              assistantMsg.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListAskHistory(userID uint, page, size int, conversationID uint) ([]response.AskHistoryItem, int64, error) {
	sessions, total, err := repository.ListAskSessionsByUser(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.AskHistoryItem, len(sessions))
	for i, s := range sessions {
		var lastQuestion, lastAnswer string
		msgs, _, _ := repository.ListAskMessages(s.ID, 1, 100)
		for _, m := range msgs {
			if m.Role == "user" {
				lastQuestion = m.Content
			} else if m.Role == "assistant" {
				lastAnswer = m.Content
			}
		}
		list[i] = response.AskHistoryItem{
			ConversationID: s.ID,
			Title:          s.Title,
			LastQuestion:   lastQuestion,
			LastAnswer:     lastAnswer,
			CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}