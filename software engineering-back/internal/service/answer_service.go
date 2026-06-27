package service

import (
	"fmt"
	"log"
	"math"
	"strings"

	"software_engineering/internal/repository"
	"software_engineering/pkg/ai"
	"software_engineering/pkg/config"
)

// AnswerService RAG 问答服务
type AnswerService struct {
	ollamaClient *ai.OllamaClient
}

// AnswerResponse 问答响应
type AnswerResponse struct {
	Answer     string           `json:"answer"`
	Confidence float64          `json:"confidence"`
	Sources    []AnswerSource   `json:"sources"`
}

// AnswerSource 参考来源
type AnswerSource struct {
	DocumentID    uint   `json:"document_id"`
	DocumentTitle string `json:"document_title"`
	Content       string `json:"content"`
}

// GraphAnswerResponse 知识图谱增强的问答响应
type GraphAnswerResponse struct {
	Answer                 string             `json:"answer"`
	Confidence             float64            `json:"confidence"`
	Sources                []AnswerSource     `json:"sources"`
	RelatedKnowledgePoints []KnowledgePoint   `json:"related_knowledge_points"`
	GraphNodesCount        int                `json:"graph_nodes_count"`
	GraphRelationsCount    int                `json:"graph_relations_count"`
}

// KnowledgePoint 知识点
type KnowledgePoint struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// answerService 问答服务单例
var answerService *AnswerService

// InitAnswerService 初始化问答服务
func InitAnswerService() {
	cfg := config.AppConfig
	ollamaClient := ai.NewOllamaClient(ai.OllamaConfig{
		BaseURL:        cfg.OllamaURL,
		Model:          cfg.OllamaModel,
		EmbeddingModel: cfg.OllamaEmbeddingModel,
	})

	answerService = &AnswerService{
		ollamaClient: ollamaClient,
	}

	log.Println("Answer service initialized")
}

// GetAnswerService 获取问答服务实例
func GetAnswerService() *AnswerService {
	return answerService
}

// SearchAndAnswer 标准 RAG 问答：向量搜索 + Rerank + LLM 回答
func (s *AnswerService) SearchAndAnswer(query string, history []ChatMessage, topK int) (*AnswerResponse, error) {
	// 1. 向量搜索（多检索一些，用于 Rerank）
	searchTopK := topK * 2
	searchResults, err := GetVectorService().Search(query, searchTopK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. Rerank 重排序
	if len(searchResults) > topK {
		searchResults, err = GetVectorService().Rerank(query, searchResults, topK)
		if err != nil {
			// Rerank 失败时使用原始排序
			if len(searchResults) > topK {
				searchResults = searchResults[:topK]
			}
		}
	}

	// 3. 构建上下文
	context := buildRAGContext(searchResults)

	// 4. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 5. 构建 prompt
	userPrompt := BuildUserPrompt(query, context, "", historyStr)

	// 6. 调用 LLM
	response, err := s.ollamaClient.Generate(userPrompt, DocumentRAGPrompt, &ai.GenerateOptions{
		Temperature: 0.7,
		TopP:        0.9,
		NumPredict:  1024,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// 7. 计算置信度
	confidence := calculateConfidence(searchResults)

	// 8. 构建来源
	sources := buildSources(searchResults)

	return &AnswerResponse{
		Answer:     response,
		Confidence: confidence,
		Sources:    sources,
	}, nil
}

// SearchAndAnswerWithGraph 知识图谱增强的 RAG 问答
func (s *AnswerService) SearchAndAnswerWithGraph(query string, history []ChatMessage, topK int) (*GraphAnswerResponse, error) {
	// 1. 向量搜索（多检索一些，用于 Rerank）
	searchTopK := topK * 2
	searchResults, err := GetVectorService().Search(query, searchTopK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. Rerank 重排序
	if len(searchResults) > topK {
		searchResults, err = GetVectorService().Rerank(query, searchResults, topK)
		if err != nil {
			// Rerank 失败时使用原始排序
			if len(searchResults) > topK {
				searchResults = searchResults[:topK]
			}
		}
	}

	// 3. 搜索知识图谱
	graphNodes, graphRelations := searchKnowledgeGraph(query)

	// 4. 构建上下文
	vectorContext := buildRAGContext(searchResults)
	graphContext := buildGraphContext(graphNodes, graphRelations)

	// 5. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 6. 构建 prompt
	combinedContext := "知识库内容：\n" + vectorContext + "\n\n知识图谱信息：\n" + graphContext
	userPrompt := BuildUserPrompt(query, combinedContext, "", historyStr)

	// 7. 调用 LLM
	response, err := s.ollamaClient.Generate(userPrompt, KnowledgeGraphPrompt, &ai.GenerateOptions{
		Temperature: 0.7,
		TopP:        0.9,
		NumPredict:  1024,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// 8. 计算置信度
	confidence := calculateConfidence(searchResults)

	// 9. 构建来源
	sources := buildSources(searchResults)

	// 10. 构建相关知识点
	relatedPoints := buildRelatedPoints(graphNodes)

	return &GraphAnswerResponse{
		Answer:                 response,
		Confidence:             confidence,
		Sources:                sources,
		RelatedKnowledgePoints: relatedPoints,
		GraphNodesCount:        len(graphNodes),
		GraphRelationsCount:    len(graphRelations),
	}, nil
}

// StreamChunk 流式响应块
type StreamChunk struct {
	Type  string `json:"type"`  // "chunk" 或 "done"
	Content string `json:"content"`
	Confidence float64 `json:"confidence,omitempty"`
}

// SearchAndAnswerStream 流式 RAG 问答
func (s *AnswerService) SearchAndAnswerStream(query string, history []ChatMessage, topK int) (<-chan StreamChunk, error) {
	// 1. 向量搜索（多检索一些，用于 Rerank）
	searchTopK := topK * 2
	searchResults, err := GetVectorService().Search(query, searchTopK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. Rerank 重排序
	if len(searchResults) > topK {
		searchResults, err = GetVectorService().Rerank(query, searchResults, topK)
		if err != nil {
			// Rerank 失败时使用原始排序
			if len(searchResults) > topK {
				searchResults = searchResults[:topK]
			}
		}
	}

	// 3. 构建上下文
	context := buildRAGContext(searchResults)

	// 4. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 5. 构建 prompt
	userPrompt := BuildUserPrompt(query, context, "", historyStr)

	// 6. 调用流式 LLM
	stream, err := s.ollamaClient.GenerateStream(userPrompt, DocumentRAGPrompt, &ai.GenerateOptions{
		Temperature: 0.7,
		TopP:        0.9,
		NumPredict:  1024,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM generation stream failed: %w", err)
	}

	// 7. 计算置信度
	confidence := calculateConfidence(searchResults)

	// 8. 创建输出 channel
	ch := make(chan StreamChunk, 100)

	go func() {
		defer close(ch)

		for chunk := range stream {
			if chunk.Done {
				// 流结束，发送完成事件
				ch <- StreamChunk{
					Type:       "done",
					Content:    "",
					Confidence: confidence,
				}
			} else {
				// 发送内容块
				content := chunk.Response
				if content == "" {
					content = chunk.Thinking
				}
				if content != "" {
					ch <- StreamChunk{
						Type:    "chunk",
						Content: content,
					}
				}
			}
		}
	}()

	return ch, nil
}

// SearchAndAnswerWithGraphStream 流式知识图谱增强 RAG 问答
func (s *AnswerService) SearchAndAnswerWithGraphStream(query string, history []ChatMessage, topK int) (<-chan StreamChunk, error) {
	// 1. 向量搜索（多检索一些，用于 Rerank）
	searchTopK := topK * 2
	searchResults, err := GetVectorService().Search(query, searchTopK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. Rerank 重排序
	if len(searchResults) > topK {
		searchResults, err = GetVectorService().Rerank(query, searchResults, topK)
		if err != nil {
			// Rerank 失败时使用原始排序
			if len(searchResults) > topK {
				searchResults = searchResults[:topK]
			}
		}
	}

	// 3. 搜索知识图谱
	graphNodes, graphRelations := searchKnowledgeGraph(query)

	// 4. 构建上下文
	vectorContext := buildRAGContext(searchResults)
	graphContext := buildGraphContext(graphNodes, graphRelations)

	// 5. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 6. 构建 prompt
	combinedContext := "知识库内容：\n" + vectorContext + "\n\n知识图谱信息：\n" + graphContext
	userPrompt := BuildUserPrompt(query, combinedContext, "", historyStr)

	// 7. 调用流式 LLM
	stream, err := s.ollamaClient.GenerateStream(userPrompt, KnowledgeGraphPrompt, &ai.GenerateOptions{
		Temperature: 0.7,
		TopP:        0.9,
		NumPredict:  1024,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM generation stream failed: %w", err)
	}

	// 8. 计算置信度
	confidence := calculateConfidence(searchResults)

	// 9. 创建输出 channel
	ch := make(chan StreamChunk, 100)

	go func() {
		defer close(ch)

		for chunk := range stream {
			if chunk.Done {
				// 流结束，发送完成事件
				ch <- StreamChunk{
					Type:       "done",
					Content:    "",
					Confidence: confidence,
				}
			} else {
				// 发送内容块
				content := chunk.Response
				if content == "" {
					content = chunk.Thinking
				}
				if content != "" {
					ch <- StreamChunk{
						Type:    "chunk",
						Content: content,
					}
				}
			}
		}
	}()

	return ch, nil
}

// buildRAGContext 构建 RAG 上下文
func buildRAGContext(results []SearchResult) string {
	if len(results) == 0 {
		return "未找到相关内容"
	}

	var sb strings.Builder
	for i, r := range results {
		sb.WriteString(fmt.Sprintf("【%d】文档ID: %d\n内容: %s\n相似度: %.2f\n\n",
			i+1, r.Metadata.DocumentID, truncateText(r.Metadata.ChunkText, 300), r.Score))
	}
	return sb.String()
}

// searchKnowledgeGraph 搜索知识图谱
func (s *AnswerService) searchKnowledgeGraphLocal(query string) ([]GraphNode, []GraphEdge) {
	// 这里简化实现，实际应该从数据库查询
	// 可以调用 graph_repo 的方法
	return []GraphNode{}, []GraphEdge{}
}

// GraphNode 图节点
type GraphNode struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// GraphEdge 图边
type GraphEdge struct {
	Source       uint   `json:"source"`
	Target       uint   `json:"target"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
}

// buildGraphContext 构建知识图谱上下文
func buildGraphContext(nodes []GraphNode, edges []GraphEdge) string {
	if len(nodes) == 0 {
		return "未找到相关知识图谱信息"
	}

	var sb strings.Builder
	sb.WriteString("相关知识点：\n")
	for _, node := range nodes {
		sb.WriteString(fmt.Sprintf("- %s (%s): %s\n", node.Name, node.Category, node.Description))
	}

	if len(edges) > 0 {
		sb.WriteString("\n知识点关系：\n")
		for _, edge := range edges {
			sb.WriteString(fmt.Sprintf("- %s --[%s]--> %s: %s\n",
				getNodeName(nodes, edge.Source), edge.RelationType,
				getNodeName(nodes, edge.Target), edge.Description))
		}
	}

	return sb.String()
}

// getNodeName 获取节点名称
func getNodeName(nodes []GraphNode, id uint) string {
	for _, node := range nodes {
		if node.ID == id {
			return node.Name
		}
	}
	return fmt.Sprintf("Node_%d", id)
}

// calculateConfidence 计算回答置信度
func calculateConfidence(results []SearchResult) float64 {
	if len(results) == 0 {
		return 0
	}

	// 计算平均相似度
	totalScore := 0.0
	for _, r := range results {
		totalScore += r.Score
	}
	avgScore := totalScore / float64(len(results))

	// 转换为置信度 (0-1)
	confidence := math.Min(avgScore*1.2, 1.0)
	return math.Round(confidence*100) / 100
}

// buildSources 构建参考来源
func buildSources(results []SearchResult) []AnswerSource {
	sources := make([]AnswerSource, 0)
	for _, r := range results {
		source := AnswerSource{
			DocumentID: r.Metadata.DocumentID,
			Content:    truncateText(r.Metadata.ChunkText, 200),
		}
		// 查询文档标题
		if r.Metadata.DocumentID > 0 {
			doc, err := repository.FindDocumentByID(r.Metadata.DocumentID)
			if err == nil {
				source.DocumentTitle = doc.Title
			}
		}
		if source.DocumentTitle == "" {
			source.DocumentTitle = "知识库"
		}
		sources = append(sources, source)
	}
	return sources
}

// buildRelatedPoints 构建相关知识点
func buildRelatedPoints(nodes []GraphNode) []KnowledgePoint {
	points := make([]KnowledgePoint, 0)
	for _, node := range nodes {
		points = append(points, KnowledgePoint{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
		})
	}
	return points
}

// truncateText 截断文本
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// searchKnowledgeGraph 搜索知识图谱的全局函数
func searchKnowledgeGraph(query string) ([]GraphNode, []GraphEdge) {
	nodes := make([]GraphNode, 0)
	edges := make([]GraphEdge, 0)

	// 1. 获取所有知识点
	allPoints, err := repository.GetAllKnowledgePoints()
	if err != nil || len(allPoints) == 0 {
		return nodes, edges
	}

	// 2. 提取关键词并匹配知识点
	queryLower := strings.ToLower(query)
	matchedIDs := make(map[uint]bool)

	for _, p := range allPoints {
		nameLower := strings.ToLower(p.Name)
		descLower := strings.ToLower(p.Description)
		if strings.Contains(nameLower, queryLower) || strings.Contains(descLower, queryLower) {
			if !matchedIDs[p.ID] {
				matchedIDs[p.ID] = true
				nodes = append(nodes, GraphNode{
					ID:          p.ID,
					Name:        p.Name,
					Description: p.Description,
					Category:    p.Category,
				})
			}
		}
	}

	// 3. 如果没有精确匹配，尝试部分关键词匹配
	if len(nodes) == 0 {
		keywords := extractKeywords(queryLower)
		for _, p := range allPoints {
			nameLower := strings.ToLower(p.Name)
			descLower := strings.ToLower(p.Description)
			for _, kw := range keywords {
				if len(kw) >= 2 && (strings.Contains(nameLower, kw) || strings.Contains(descLower, kw)) {
					if !matchedIDs[p.ID] {
						matchedIDs[p.ID] = true
						nodes = append(nodes, GraphNode{
							ID:          p.ID,
							Name:        p.Name,
							Description: p.Description,
							Category:    p.Category,
						})
					}
					break
				}
			}
		}
	}

	// 4. 获取匹配知识点的关系
	if len(matchedIDs) > 0 {
		_, allRelations, _ := repository.GetAllGraphDataFromNeo4j()
		for _, rel := range allRelations {
			if matchedIDs[rel.SourceID] || matchedIDs[rel.TargetID] {
				edges = append(edges, GraphEdge{
					Source:       rel.SourceID,
					Target:       rel.TargetID,
					RelationType: rel.RelationType,
					Description:  rel.Description,
				})
			}
		}
	}

	return nodes, edges
}

// Ensure AnswerService uses ai package
var _ = (*ai.OllamaClient)(nil)
