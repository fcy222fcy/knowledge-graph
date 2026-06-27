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

// CreateSession 创建新的问答会话
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

// ListSessions 分页获取用户的问答会话列表，包含最后一条消息摘要
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

// ListSessionMessages 分页获取指定会话的消息记录
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

// Ask 智能问答核心方法，采用多级降级策略：
// 1. 知识图谱查询（Graph RAG）→ 2. 向量检索 RAG → 3. 关键词检索 → 4. 知识点匹配 → 5. LLM 自由回答
func Ask(userID uint, req request.AskRequest) (*response.AskResponse, error) {
	log.Printf("DEBUG: Received question: %q", req.Question)

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

	// 获取历史消息用于上下文（只取最近2条用户问题，不包含AI回答避免干扰）
	historyMsgs, _ := repository.ListRecentMessages(sessionID, 10)
	history := make([]ChatMessage, 0)
	for _, m := range historyMsgs {
		// 只保留用户消息作为上下文，不包含AI回答
		if m.Role == "user" {
			history = append(history, ChatMessage{Role: m.Role, Content: m.Content})
		}
	}
	// 只保留最近2条历史
	if len(history) > 2 {
		history = history[len(history)-2:]
	}

	// 保存当前用户消息（防止重复保存）
	userMsg := &entity.AskMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Question,
	}
	if !repository.HasRecentUserMessage(sessionID, req.Question, 5) {
		repository.CreateAskMessage(userMsg)
	}

	// 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	var answer string
	var confidence float64
	var sources []response.AskSource
	var related []response.KPRef

	// 第1级：优先使用知识图谱查询
	graphContext, graphSources, graphRelated, graphConfidence := searchFromKnowledgeGraph(req.Question)
	if graphContext != "" {
		log.Printf("DEBUG: Using knowledge graph for question: %s", req.Question)
		// 尝试调用 LLM 基于图谱上下文生成回答
		if aiClient.IsAvailable() {
			userPrompt := BuildGraphUserPrompt(req.Question, graphContext, historyStr)

			llmResponse, err := aiClient.Generate(userPrompt, KnowledgeGraphPrompt, nil)
			if err == nil && llmResponse != "" {
				answer = llmResponse
				confidence = graphConfidence
				sources = graphSources
				related = graphRelated
			}
		}

		// LLM 不可用或失败时，直接返回图谱上下文
		if answer == "" {
			answer = fmt.Sprintf("关于「%s」的知识图谱查询结果：\n\n%s\n\n以上内容来自知识图谱。", req.Question, graphContext)
			confidence = 0.75
			sources = graphSources
			related = graphRelated
		}
	}

	// 第2级：向量检索 RAG
	vecSvc := GetVectorService()
	if answer == "" && vecSvc != nil && vecSvc.GetSize() > 0 {
		log.Printf("DEBUG: Using vector search for question: %s", req.Question)
		searchResults, err := vecSvc.Search(req.Question, 5)
		if err == nil && len(searchResults) > 0 {
			// 过滤低相关性结果（相似度 < 0.5）
			var validResults []SearchResult
			for _, r := range searchResults {
				if r.Score >= 0.5 {
					validResults = append(validResults, r)
				}
			}

			if len(validResults) > 0 {
				// 取 top-3 相关片段
				contextParts := make([]string, 0, 3)
				for i, r := range validResults {
					if i >= 3 {
						break
					}
					contextParts = append(contextParts, r.Metadata.ChunkText)
				}
				contextText := strings.Join(contextParts, "\n\n---\n\n")

				// 获取文档标题
				docTitle := "知识库"
				if validResults[0].Metadata.DocumentID > 0 {
					doc, err := repository.FindDocumentByID(validResults[0].Metadata.DocumentID)
					if err == nil {
						docTitle = doc.Title
						sources = append(sources, response.AskSource{
							DocumentID:    doc.ID,
							DocumentTitle: doc.Title,
							Content:       contextText,
						})
					}
				}

				if aiClient.IsAvailable() {
					userPrompt := BuildUserPrompt(req.Question, contextText, docTitle, historyStr)

					llmResponse, err := aiClient.Generate(userPrompt, DocumentRAGPrompt, nil)
					if err == nil && llmResponse != "" {
						answer = llmResponse
						confidence = 0.9
					}
				}

				if answer == "" {
					answer = fmt.Sprintf("关于「%s」的回答：\n\n根据文档《%s》中的内容：\n\n%s\n\n以上内容来自知识库文档检索。", req.Question, docTitle, contextText)
					confidence = 0.75
				}
			}
		}
	}

	// 降级到本地关键词检索
	if answer == "" {
		log.Printf("DEBUG: Falling back to keyword search for question: %s", req.Question)
		docs, _ := repository.GetAllDocumentsContent()
		log.Printf("DEBUG: Found %d documents for keyword search", len(docs))
		questionLower := strings.ToLower(req.Question)

		// 提取关键词
		keywords := extractKeywords(questionLower)
		// 同时用原始问题作为搜索词（去掉问号等）
		cleanQuestion := strings.NewReplacer("？", "", "!", "", "？", "", "。", "", "，", "", ",", "", ".", "").Replace(questionLower)

		type docMatch struct {
			doc     entity.Document
			score   int
			snippet string
		}
		matches := make([]docMatch, 0)

		for _, doc := range docs {
			contentLower := strings.ToLower(doc.Content)
			titleLower := strings.ToLower(doc.Title)

			score := 0

			// 1. 标题匹配（最高优先级）
			for _, kw := range keywords {
				if len(kw) >= 2 && strings.Contains(titleLower, kw) {
					score += 5
				}
			}

			// 2. 完整问题匹配（高分）
			if strings.Contains(contentLower, cleanQuestion) {
				score += 10
			}

			// 3. 关键词逐个匹配
			for _, kw := range keywords {
				if len(kw) >= 2 && strings.Contains(contentLower, kw) {
					score++
				}
			}

			// 3. 如果关键词太长没有匹配，尝试拆分成2字子串匹配
			if score == 0 {
				for _, kw := range keywords {
					if len(kw) > 2 {
						// 将关键词拆成所有可能的2字及以上的子串
						for segLen := 2; segLen <= len(kw); segLen++ {
							for startIdx := 0; startIdx <= len(kw)-segLen; startIdx++ {
								sub := kw[startIdx : startIdx+segLen]
								if strings.Contains(contentLower, sub) {
									score++
									break
								}
							}
							if score > 0 {
								break
							}
						}
					}
				}
			}

			if score > 0 {
				// 提取匹配位置的上下文，按 markdown 标题智能截取整个 section
				idx := -1
				for _, kw := range keywords {
					if i := strings.Index(contentLower, kw); i >= 0 {
						idx = i
						break
					}
				}
				if idx < 0 {
					idx = 0
				}
				snippet := extractSectionByHeader(doc.Content, idx)
				matches = append(matches, docMatch{doc: doc, score: score, snippet: snippet})
			}
		}

		log.Printf("DEBUG: Found %d document matches for question: %s", len(matches), req.Question)

		// 过滤低分匹配（至少需要匹配1个关键词）
		var validMatches []docMatch
		for _, m := range matches {
			log.Printf("DEBUG: Document %d (%s) score: %d", m.doc.ID, m.doc.Title, m.score)
			if m.score >= 1 {
				validMatches = append(validMatches, m)
			}
		}
		matches = validMatches

		// 按匹配分数排序，融合多个片段
		if len(matches) > 0 {
			// 排序取 top-3
			for i := 0; i < len(matches)-1; i++ {
				for j := i + 1; j < len(matches); j++ {
					if matches[j].score > matches[i].score {
						matches[i], matches[j] = matches[j], matches[i]
					}
				}
			}
			topN := min(3, len(matches))

			// 融合多个片段
			contextParts := make([]string, 0, topN)
			docTitle := matches[0].doc.Title
			docID := matches[0].doc.ID
			for i := 0; i < topN; i++ {
				contextParts = append(contextParts, matches[i].snippet)
			}
			contextText := strings.Join(contextParts, "\n\n---\n\n")

			sources = append(sources, response.AskSource{
				DocumentID:    docID,
				DocumentTitle: docTitle,
				Content:       contextText,
			})

			// 尝试调用 LLM 生成回答
			if aiClient.IsAvailable() {
				userPrompt := BuildUserPrompt(req.Question, contextText, docTitle, historyStr)

				llmResponse, err := aiClient.Generate(userPrompt, DocumentRAGPrompt, nil)
				if err == nil && llmResponse != "" {
					answer = llmResponse
					confidence = 0.85
				} else {
					// LLM 失败，降级为直接返回文档片段
					log.Printf("warning: LLM generation failed for local search: %v", err)
					answer = fmt.Sprintf("关于「%s」的回答：\n\n根据文档《%s》中的内容：\n\n%s\n\n📚 **参考来源**：《%s》", req.Question, docTitle, contextText, docTitle)
					confidence = 0.7
				}
			} else {
				// AI 服务不可用，直接返回文档片段
				answer = fmt.Sprintf("关于「%s」的回答：\n\n根据文档《%s》中的内容：\n\n%s\n\n📚 **参考来源**：《%s》", req.Question, docTitle, contextText, docTitle)
				confidence = 0.7
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
		}
	}

	// 最终降级：如果知识库没有找到，让 LLM 自由回答
	if answer == "" && aiClient.IsAvailable() {
		log.Printf("info: No knowledge base match for question: %s, using LLM free answer", req.Question)
		userPrompt := BuildFreeUserPrompt(req.Question)

		llmResponse, err := aiClient.Generate(userPrompt, FreeAnswerPrompt, nil)
		if err == nil && llmResponse != "" {
			answer = llmResponse
			confidence = 0.5
		}
	}

	// 最终兜底
	if answer == "" {
		answer = fmt.Sprintf("抱歉，暂时无法回答关于「%s」的问题。您可以尝试：\n1. 上传更多相关文档\n2. 构建知识图谱\n3. 联系管理员获取帮助", req.Question)
		confidence = 0.3
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

// ListAskHistory 获取用户的问答历史记录列表
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

// extractSectionByHeader 按 markdown 标题层级智能截取内容片段。
// 从匹配位置向前找到当前所属的标题，向后截取到同级或更高级标题出现之前。
func extractSectionByHeader(content string, matchIdx int) string {
	if matchIdx < 0 || matchIdx >= len(content) {
		return content
	}

	lines := strings.Split(content, "\n")

	// 找到 matchIdx 所在的行号
	charCount := 0
	matchLine := 0
	for i, line := range lines {
		charCount += len(line) + 1 // +1 for \n
		if charCount > matchIdx {
			matchLine = i
			break
		}
	}

	// 从匹配位置向上找当前所属标题
	type headerInfo struct {
		level int
		index int
	}
	currentHeader := headerInfo{-1, 0}
	for i := matchLine; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "#") {
			level := 0
			for _, ch := range line {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			if level > 0 {
				currentHeader = headerInfo{level, i}
				break
			}
		}
	}

	// 从当前标题的下一行开始，向后找到同级或更高级标题
	endLine := len(lines)
	if currentHeader.level > 0 {
		for i := currentHeader.index + 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if strings.HasPrefix(line, "#") {
				level := 0
				for _, ch := range line {
					if ch == '#' {
						level++
					} else {
						break
					}
				}
				// 遇到同级或更高级标题，截断
				if level > 0 && level <= currentHeader.level {
					endLine = i
					break
				}
			}
		}
	} else {
		// 没找到标题，回退到原始逻辑：前后各取一定字符
		startLine := max(0, matchLine-3)
		endLine = min(len(lines), matchLine+10)
		return strings.Join(lines[startLine:endLine], "\n")
	}

	section := strings.Join(lines[currentHeader.index:endLine], "\n")
	section = strings.TrimSpace(section)

	// 安全上限：防止超大 section，返回精准片段
	if len(section) > 1500 {
		section = section[:1500]
	}
	return section
}

// extractKeywords 从问题中提取关键词
func extractKeywords(question string) []string {
	// 常见停用词
	stopWords := map[string]bool{
		"的": true, "了": true, "在": true, "是": true, "我": true,
		"有": true, "和": true, "就": true, "不": true, "人": true,
		"都": true, "一": true, "一个": true, "上": true, "也": true,
		"很": true, "到": true, "说": true, "要": true, "去": true,
		"你": true, "会": true, "着": true, "没有": true, "看": true,
		"好": true, "自己": true, "这": true, "他": true, "她": true,
		"它": true, "们": true, "那": true, "怎么": true, "什么": true,
		"如何": true, "怎样": true, "可以": true, "能": true, "请": true,
		"帮": true, "帮我": true, "告诉": true, "一下": true, "下": true,
		"吗": true, "呢": true, "吧": true, "啊": true,
	}

	// 移除常见问句模式
	question = strings.ReplaceAll(question, "怎么写", "")
	question = strings.ReplaceAll(question, "如何写", "")
	question = strings.ReplaceAll(question, "怎么用", "")
	question = strings.ReplaceAll(question, "如何使用", "")
	question = strings.ReplaceAll(question, "是什么", "")
	question = strings.ReplaceAll(question, "有哪些", "")
	question = strings.ReplaceAll(question, "怎么", "")
	question = strings.ReplaceAll(question, "如何", "")

	// 按空格和标点分割
	words := strings.FieldsFunc(question, func(r rune) bool {
		return r == ' ' || r == '，' || r == '。' || r == '？' || r == '！' ||
			r == ',' || r == '.' || r == '?' || r == '!' || r == '、'
	})

	keywords := make([]string, 0)
	seen := make(map[string]bool)

	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) < 2 {
			continue
		}
		if stopWords[word] {
			continue
		}
		if !seen[word] {
			seen[word] = true
			keywords = append(keywords, word)
		}
	}

	// 如果提取不到关键词，返回原始问题的子串
	if len(keywords) == 0 && len(question) >= 2 {
		keywords = append(keywords, question)
	}

	return keywords
}

// searchFromKnowledgeGraph 从知识图谱查询与问题相关的知识点和关系，构建上下文
// 返回：图谱上下文文本、参考来源、相关知识点列表、置信度
func searchFromKnowledgeGraph(question string) (string, []response.AskSource, []response.KPRef, float64) {
	// 从 Neo4j 获取图谱数据（不可用时自动降级到 MySQL）
	allNodes, allRelations, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil || len(allNodes) == 0 {
		log.Printf("DEBUG: Knowledge graph query returned no nodes or error: %v", err)
		return "", nil, nil, 0
	}

	// 提取关键词并匹配知识点
	keywords := extractKeywords(strings.ToLower(question))
	if len(keywords) == 0 {
		cleanQuestion := strings.NewReplacer("？", "", "!", "", "。", "", "，", "", ",", "", ".", "").Replace(strings.ToLower(question))
		if len(cleanQuestion) >= 2 {
			keywords = append(keywords, cleanQuestion)
		}
	}

	// 根据关键词匹配知识点
	matchedNodeIDs := make(map[uint]bool)
	var matchedNodes []entity.KnowledgePoint
	for _, node := range allNodes {
		nameLower := strings.ToLower(node.Name)
		descLower := strings.ToLower(node.Description)
		for _, kw := range keywords {
			if len(kw) >= 2 && (strings.Contains(nameLower, kw) || strings.Contains(descLower, kw)) {
				if !matchedNodeIDs[node.ID] {
					matchedNodeIDs[node.ID] = true
					matchedNodes = append(matchedNodes, node)
				}
				break
			}
		}
	}

	if len(matchedNodes) == 0 {
		log.Printf("DEBUG: No knowledge graph nodes matched for question: %s", question)
		return "", nil, nil, 0
	}

	// 获取与匹配节点相关的边
	var matchedRelations []entity.KnowledgeRelation
	for _, rel := range allRelations {
		if matchedNodeIDs[rel.SourceID] || matchedNodeIDs[rel.TargetID] {
			matchedRelations = append(matchedRelations, rel)
		}
	}

	// 构建图谱上下文
	var sb strings.Builder
	sb.WriteString("相关知识点：\n")
	for _, node := range matchedNodes {
		category := node.Category
		if category == "" {
			category = "未分类"
		}
		sb.WriteString(fmt.Sprintf("- %s (%s): %s\n", node.Name, category, node.Description))
	}

	if len(matchedRelations) > 0 {
		sb.WriteString("\n知识点关系：\n")
		nodeIDToName := make(map[uint]string)
		for _, node := range allNodes {
			nodeIDToName[node.ID] = node.Name
		}
		for _, rel := range matchedRelations {
			sourceName := nodeIDToName[rel.SourceID]
			if sourceName == "" {
				sourceName = fmt.Sprintf("节点_%d", rel.SourceID)
			}
			targetName := nodeIDToName[rel.TargetID]
			if targetName == "" {
				targetName = fmt.Sprintf("节点_%d", rel.TargetID)
			}
			relType := rel.RelationType
			if relType == "" {
				relType = rel.Type
			}
			sb.WriteString(fmt.Sprintf("- %s --[%s]--> %s: %s\n", sourceName, relType, targetName, rel.Description))
		}
	}

	graphContext := sb.String()

	// 构建 sources
	sources := []response.AskSource{
		{
			DocumentTitle: "知识图谱",
			Content:       graphContext,
		},
	}

	// 构建相关知识点
	related := make([]response.KPRef, 0, len(matchedNodes))
	for _, node := range matchedNodes {
		related = append(related, response.KPRef{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
		})
	}

	log.Printf("DEBUG: Knowledge graph found %d nodes and %d relations", len(matchedNodes), len(matchedRelations))
	return graphContext, sources, related, 0.85
}

// StreamEvent 流式事件
type StreamEvent struct {
	Type       string             `json:"type"`       // "chunk", "done", "error"
	Content    string             `json:"content"`
	Confidence float64            `json:"confidence,omitempty"`
	Sources    []response.AskSource `json:"sources,omitempty"`
	Related    []response.KPRef   `json:"related,omitempty"`
}

// AskStream 流式智能问答核心方法
func AskStream(userID uint, req request.AskRequest) (uint, <-chan StreamEvent, error) {
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

	// 获取历史消息用于上下文（只取最近2条用户问题，不包含AI回答避免干扰）
	historyMsgs, _ := repository.ListRecentMessages(sessionID, 10)
	history := make([]ChatMessage, 0)
	for _, m := range historyMsgs {
		// 只保留用户消息作为上下文，不包含AI回答
		if m.Role == "user" {
			history = append(history, ChatMessage{Role: m.Role, Content: m.Content})
		}
	}
	// 只保留最近2条历史
	if len(history) > 2 {
		history = history[len(history)-2:]
	}

	// 保存当前用户消息（防止重复保存）
	if !repository.HasRecentUserMessage(sessionID, req.Question, 5) {
		userMsg := &entity.AskMessage{
			SessionID: sessionID,
			Role:      "user",
			Content:   req.Question,
		}
		repository.CreateAskMessage(userMsg)
	}

	// 创建事件 channel
	ch := make(chan StreamEvent, 100)

	go func() {
		defer close(ch)

		var answer strings.Builder
		var confidence float64
		var sources []response.AskSource
		var related []response.KPRef

		if aiClient.IsAvailable() {
			// 尝试使用知识图谱流式问答
			stream, err := aiClient.SearchAndAnswerWithGraphStream(req.Question, history, 3)
			if err == nil {
				for chunk := range stream {
					if chunk.Type == "done" {
						confidence = chunk.Confidence
					} else if chunk.Type == "chunk" {
						answer.WriteString(chunk.Content)
						ch <- StreamEvent{
							Type:    "chunk",
							Content: chunk.Content,
						}
					}
				}
			} else {
				// 降级到普通 RAG 流式问答
				log.Printf("warning: Graph stream QA failed, falling back to RAG stream: %v", err)
				stream, err := aiClient.SearchAndAnswerWithHistoryStream(req.Question, history, 3)
				if err == nil {
					for chunk := range stream {
						if chunk.Type == "done" {
							confidence = chunk.Confidence
						} else if chunk.Type == "chunk" {
							answer.WriteString(chunk.Content)
							ch <- StreamEvent{
								Type:    "chunk",
								Content: chunk.Content,
							}
						}
					}
				} else {
					// 降级到非流式问答
					log.Printf("warning: Stream QA failed, falling back to non-stream: %v", err)
					resp, err := Ask(userID, req)
					if err == nil {
						ch <- StreamEvent{
							Type:    "chunk",
							Content: resp.Answer,
						}
						confidence = resp.Confidence
						sources = resp.Sources
						related = resp.RelatedKnowledgePoints
					} else {
						ch <- StreamEvent{
							Type:    "error",
							Content: "问答服务暂时不可用，请稍后重试",
						}
					}
				}
			}
		} else {
			// AI 服务不可用，使用非流式降级
			resp, err := Ask(userID, req)
			if err == nil {
				ch <- StreamEvent{
					Type:    "chunk",
					Content: resp.Answer,
				}
				confidence = resp.Confidence
				sources = resp.Sources
				related = resp.RelatedKnowledgePoints
			} else {
				ch <- StreamEvent{
					Type:    "error",
					Content: "问答服务暂时不可用，请稍后重试",
				}
			}
		}

		// 保存助手消息
		assistantMsg := &entity.AskMessage{
			SessionID:  sessionID,
			Role:       "assistant",
			Content:    answer.String(),
			Confidence: confidence,
		}
		repository.CreateAskMessage(assistantMsg)

		// 发送完成事件
		ch <- StreamEvent{
			Type:       "done",
			Content:    "",
			Confidence: confidence,
			Sources:    sources,
			Related:    related,
		}
	}()

	return sessionID, ch, nil
}