package service

import (
	"fmt"
	"log"
	"math"
	"strings"

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

// SearchAndAnswer 标准 RAG 问答：向量搜索 + LLM 回答
func (s *AnswerService) SearchAndAnswer(query string, history []ChatMessage, topK int) (*AnswerResponse, error) {
	// 1. 向量搜索
	searchResults, err := GetVectorService().Search(query, topK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. 构建上下文
	context := buildRAGContext(searchResults)

	// 3. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 4. 构建 prompt
	systemPrompt := `你是一个软件工程知识问答助手。请根据提供的知识库内容回答用户的问题。
如果知识库中没有相关信息，请说明无法从知识库中找到答案。
回答要准确、简洁，并引用相关的知识来源。`

	userPrompt := fmt.Sprintf(`知识库内容：
%s

%s
用户问题：%s

请基于以上知识库内容回答问题：`, context, historyStr, query)

	// 5. 调用 LLM
	response, err := s.ollamaClient.Generate(userPrompt, systemPrompt, &ai.GenerateOptions{
		Temperature: 0.7,
		TopP:        0.9,
		NumPredict:  1024,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// 6. 计算置信度
	confidence := calculateConfidence(searchResults)

	// 7. 构建来源
	sources := buildSources(searchResults)

	return &AnswerResponse{
		Answer:     response,
		Confidence: confidence,
		Sources:    sources,
	}, nil
}

// SearchAndAnswerWithGraph 知识图谱增强的 RAG 问答
func (s *AnswerService) SearchAndAnswerWithGraph(query string, history []ChatMessage, topK int) (*GraphAnswerResponse, error) {
	// 1. 向量搜索
	searchResults, err := GetVectorService().Search(query, topK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. 搜索知识图谱
	graphNodes, graphRelations := searchKnowledgeGraph(query)

	// 3. 构建上下文
	vectorContext := buildRAGContext(searchResults)
	graphContext := buildGraphContext(graphNodes, graphRelations)

	// 4. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 5. 构建 prompt
	systemPrompt := `你是一个软件工程知识问答助手。请根据提供的知识库内容和知识图谱信息回答用户的问题。
知识图谱提供了知识点之间的关系，可以帮助你更好地理解概念之间的联系。
如果知识库中没有相关信息，请说明无法从知识库中找到答案。
回答要准确、简洁，并引用相关的知识来源。`

	userPrompt := fmt.Sprintf(`知识库内容：
%s

知识图谱信息：
%s

%s
用户问题：%s

请基于以上信息回答问题：`, vectorContext, graphContext, historyStr, query)

	// 6. 调用 LLM
	response, err := s.ollamaClient.Generate(userPrompt, systemPrompt, &ai.GenerateOptions{
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

	// 9. 构建相关知识点
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
	// 1. 向量搜索
	searchResults, err := GetVectorService().Search(query, topK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. 构建上下文
	context := buildRAGContext(searchResults)

	// 3. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 4. 构建 prompt
	systemPrompt := `你是一个软件工程知识问答助手。请根据提供的知识库内容回答用户的问题。
如果知识库中没有相关信息，请说明无法从知识库中找到答案。
回答要准确、简洁，并引用相关的知识来源。`

	userPrompt := fmt.Sprintf(`知识库内容：
%s

%s
用户问题：%s

请基于以上知识库内容回答问题：`, context, historyStr, query)

	// 5. 调用流式 LLM
	stream, err := s.ollamaClient.GenerateStream(userPrompt, systemPrompt, &ai.GenerateOptions{
		Temperature: 0.7,
		TopP:        0.9,
		NumPredict:  1024,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM generation stream failed: %w", err)
	}

	// 6. 计算置信度
	confidence := calculateConfidence(searchResults)

	// 7. 创建输出 channel
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
	// 1. 向量搜索
	searchResults, err := GetVectorService().Search(query, topK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 2. 搜索知识图谱
	graphNodes, graphRelations := searchKnowledgeGraph(query)

	// 3. 构建上下文
	vectorContext := buildRAGContext(searchResults)
	graphContext := buildGraphContext(graphNodes, graphRelations)

	// 4. 构建对话历史字符串
	historyStr := BuildConversationContext(history)

	// 5. 构建 prompt
	systemPrompt := `你是一个软件工程知识问答助手。请根据提供的知识库内容和知识图谱信息回答用户的问题。
知识图谱提供了知识点之间的关系，可以帮助你更好地理解概念之间的联系。
如果知识库中没有相关信息，请说明无法从知识库中找到答案。
回答要准确、简洁，并引用相关的知识来源。`

	userPrompt := fmt.Sprintf(`知识库内容：
%s

知识图谱信息：
%s

%s
用户问题：%s

请基于以上信息回答问题：`, vectorContext, graphContext, historyStr, query)

	// 6. 调用流式 LLM
	stream, err := s.ollamaClient.GenerateStream(userPrompt, systemPrompt, &ai.GenerateOptions{
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
		sources = append(sources, AnswerSource{
			DocumentID: r.Metadata.DocumentID,
			Content:    truncateText(r.Metadata.ChunkText, 200),
		})
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
	// 简化实现：从 MySQL 查询相关知识点
	// 实际应该调用 graph_repo 的方法
	nodes := make([]GraphNode, 0)
	edges := make([]GraphEdge, 0)

	// TODO: 实现从数据库查询知识图谱
	// 可以调用 repository.FindKnowledgePointsByKeyword 等方法

	return nodes, edges
}

// Ensure AnswerService uses ai package
var _ = (*ai.OllamaClient)(nil)
