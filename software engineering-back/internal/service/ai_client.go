package service

import (
	"fmt"
	"log"
	"strings"

	"software_engineering/internal/repository"
)

// AIClient AI 服务客户端，负责调用知识图谱构建、语义搜索和智能问答接口
// 现在直接调用本地 AI 服务，不再通过 HTTP 调用 Python 服务
type AIClient struct {
	available bool
}

// AIAnswerRequest AI 问答请求
type AIAnswerRequest struct {
	Query string `json:"query"` // 用户查询
	TopK  int    `json:"top_k"` // 返回结果数量
}

// AIAnswerSource AI 问答参考来源
type AIAnswerSource struct {
	DocumentID    int    `json:"document_id"`    // 文档ID
	DocumentTitle string `json:"document_title"` // 文档标题
	Content       string `json:"content"`        // 引用内容
}

// AIAnswerResponse AI 问答响应
type AIAnswerResponse struct {
	Answer     string           `json:"answer"`      // 生成的回答
	Confidence float64          `json:"confidence"`  // 置信度
	Sources    []AIAnswerSource `json:"sources"`     // 参考来源
}

// AIKnowledgePoint 相关知识点
type AIKnowledgePoint struct {
	ID          uint   `json:"id"`          // 知识点ID
	Name        string `json:"name"`        // 知识点名称
	Description string `json:"description"` // 知识点描述
}

// AIAnswerWithGraphResponse 基于知识图谱的回答响应
type AIAnswerWithGraphResponse struct {
	Answer                  string             `json:"answer"`                   // 生成的回答
	Confidence              float64            `json:"confidence"`               // 置信度
	Sources                 []AIAnswerSource   `json:"sources"`                  // 参考来源
	RelatedKnowledgePoints  []AIKnowledgePoint `json:"related_knowledge_points"` // 相关知识点
	GraphNodesCount         int                `json:"graph_nodes_count"`        // 图谱节点数
	GraphRelationsCount     int                `json:"graph_relations_count"`    // 图谱关系数
}

// ChatMessage 对话消息，用于传递历史上下文
type ChatMessage struct {
	Role    string `json:"role"`    // 消息角色：user / assistant
	Content string `json:"content"` // 消息内容
}

// AIAnswerWithHistoryRequest 带历史上下文的请求
type AIAnswerWithHistoryRequest struct {
	Query   string        `json:"query"`   // 用户查询
	History []ChatMessage `json:"history"` // 对话历史
	TopK    int           `json:"top_k"`   // 返回结果数量
}

// AIBuildRequest AI 知识图谱构建请求
type AIBuildRequest struct {
	DocumentID uint   `json:"document_id"` // 文档ID
	Title      string `json:"title"`       // 文档标题
	Content    string `json:"content"`     // 文档内容
	Source     string `json:"source"`      // 来源标识
}

// AIBuildResponse AI 知识图谱构建响应
type AIBuildResponse struct {
	DocumentID       uint          `json:"document_id"`       // 文档ID
	DocumentTitle    string        `json:"document_title"`    // 文档标题
	CreatedPoints    int           `json:"created_points"`    // 创建的知识点数
	CreatedRelations int           `json:"created_relations"` // 创建的关系数
	ChunkCount       int           `json:"chunk_count"`       // 文档分块数
	VectorCount      int           `json:"vector_count"`      // 向量数
	Status           string        `json:"status"`            // 构建状态
	Message          string        `json:"message"`           // 构建结果描述
	Points           []AIGraphNode `json:"points"`            // 知识点列表
	Relations        []AIGraphEdge `json:"relations"`         // 关系列表
}

// AISearchRequest AI 语义搜索请求
type AISearchRequest struct {
	Query string `json:"query"` // 搜索查询
	TopK  int    `json:"top_k"` // 返回结果数量
}

// AISearchResult AI 语义搜索结果
type AISearchResult struct {
	ChunkText         string  `json:"chunk_text"`          // 文档片段文本
	Score             float64 `json:"score"`               // 相似度分数
	DocumentID        int     `json:"document_id"`         // 文档ID
	KnowledgePointIDs []int   `json:"knowledge_point_ids"` // 关联的知识点ID列表
}

// AISearchResponse AI 语义搜索响应
type AISearchResponse struct {
	Results []AISearchResult `json:"results"` // 搜索结果列表
}

// AIGraphResponse AI 知识图谱响应
type AIGraphResponse struct {
	Nodes []AIGraphNode `json:"nodes"` // 节点列表
	Edges []AIGraphEdge `json:"edges"` // 边列表
}

// AIGraphNode AI 知识图谱节点
type AIGraphNode struct {
	ID          uint   `json:"id"`          // 节点ID
	Name        string `json:"name"`        // 节点名称
	Description string `json:"description"` // 节点描述
	Category    string `json:"category"`    // 节点分类
	DocumentID  uint   `json:"document_id"` // 关联的文档ID
}

// AIGraphEdge AI 知识图谱边（关系）
type AIGraphEdge struct {
	Source       uint   `json:"source"`       // 源节点ID
	Target       uint   `json:"target"`       // 目标节点ID
	RelationType string `json:"relation_type"` // 关系类型
	Description  string `json:"description"`  // 关系描述
}

// IsAvailable 检查 AI 服务是否已配置
func (c *AIClient) IsAvailable() bool {
	return c.available
}

// BuildGraph 调用 AI 服务构建知识图谱
func (c *AIClient) BuildGraph(req AIBuildRequest) (*AIBuildResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not available")
	}

	// 使用抽取服务提取知识点和关系
	extractionService := GetExtractionService()
	if extractionService == nil {
		return nil, fmt.Errorf("extraction service not initialized")
	}

	result, err := extractionService.ExtractKnowledgePoints(req.Content, req.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	// 转换为 AIGraphNode 和 AIGraphEdge
	nodes := make([]AIGraphNode, 0, len(result.Points))
	pointIDMap := make(map[string]uint)
	for i, point := range result.Points {
		node := AIGraphNode{
			ID:          uint(i + 1),
			Name:        point.Name,
			Description: point.Description,
			Category:    point.Category,
			DocumentID:  req.DocumentID,
		}
		nodes = append(nodes, node)
		pointIDMap[point.Name] = node.ID
	}

	edges := make([]AIGraphEdge, 0, len(result.Relations))
	for _, rel := range result.Relations {
		sourceID, ok := pointIDMap[rel.Source]
		if !ok {
			continue
		}
		targetID, ok := pointIDMap[rel.Target]
		if !ok {
			continue
		}
		edge := AIGraphEdge{
			Source:       sourceID,
			Target:       targetID,
			RelationType: rel.RelationType,
			Description:  rel.Description,
		}
		edges = append(edges, edge)
	}

	return &AIBuildResponse{
		DocumentID:       req.DocumentID,
		DocumentTitle:    req.Title,
		CreatedPoints:    len(nodes),
		CreatedRelations: len(edges),
		Status:           "completed",
		Message:          fmt.Sprintf("成功提取 %d 个知识点和 %d 个关系", len(nodes), len(edges)),
		Points:           nodes,
		Relations:        edges,
	}, nil
}

// Search 调用 AI 服务进行语义搜索
func (c *AIClient) Search(query string, topK int) (*AISearchResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not available")
	}

	vectorService := GetVectorService()
	if vectorService == nil {
		return nil, fmt.Errorf("vector service not initialized")
	}

	results, err := vectorService.Search(query, topK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	searchResults := make([]AISearchResult, 0, len(results))
	for _, r := range results {
		searchResults = append(searchResults, AISearchResult{
			ChunkText:  r.Metadata.ChunkText,
			Score:      r.Score,
			DocumentID: int(r.Metadata.DocumentID),
		})
	}

	return &AISearchResponse{Results: searchResults}, nil
}

// SearchAndAnswer 调用 AI 服务进行语义搜索并生成回答
func (c *AIClient) SearchAndAnswer(query string, topK int) (*AIAnswerResponse, error) {
	return c.SearchAndAnswerWithHistory(query, nil, topK)
}

// SearchAndAnswerWithHistory 带对话历史的智能问答
func (c *AIClient) SearchAndAnswerWithHistory(query string, history []ChatMessage, topK int) (*AIAnswerResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not available")
	}

	answerService := GetAnswerService()
	if answerService == nil {
		return nil, fmt.Errorf("answer service not initialized")
	}

	resp, err := answerService.SearchAndAnswer(query, history, topK)
	if err != nil {
		return nil, err
	}

	sources := make([]AIAnswerSource, 0, len(resp.Sources))
	for _, s := range resp.Sources {
		sources = append(sources, AIAnswerSource{
			DocumentID: int(s.DocumentID),
			Content:    s.Content,
		})
	}

	return &AIAnswerResponse{
		Answer:     resp.Answer,
		Confidence: resp.Confidence,
		Sources:    sources,
	}, nil
}

// SearchAndAnswerWithGraph 基于知识图谱的智能问答
func (c *AIClient) SearchAndAnswerWithGraph(query string, history []ChatMessage, topK int) (*AIAnswerWithGraphResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not available")
	}

	answerService := GetAnswerService()
	if answerService == nil {
		return nil, fmt.Errorf("answer service not initialized")
	}

	resp, err := answerService.SearchAndAnswerWithGraph(query, history, topK)
	if err != nil {
		return nil, err
	}

	sources := make([]AIAnswerSource, 0, len(resp.Sources))
	for _, s := range resp.Sources {
		sources = append(sources, AIAnswerSource{
			DocumentID: int(s.DocumentID),
			Content:    s.Content,
		})
	}

	points := make([]AIKnowledgePoint, 0, len(resp.RelatedKnowledgePoints))
	for _, p := range resp.RelatedKnowledgePoints {
		points = append(points, AIKnowledgePoint{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return &AIAnswerWithGraphResponse{
		Answer:                 resp.Answer,
		Confidence:             resp.Confidence,
		Sources:                sources,
		RelatedKnowledgePoints: points,
		GraphNodesCount:        resp.GraphNodesCount,
		GraphRelationsCount:    resp.GraphRelationsCount,
	}, nil
}

// SearchAndAnswerWithHistoryStream 带对话历史的流式智能问答
func (c *AIClient) SearchAndAnswerWithHistoryStream(query string, history []ChatMessage, topK int) (<-chan StreamChunk, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not available")
	}

	answerService := GetAnswerService()
	if answerService == nil {
		return nil, fmt.Errorf("answer service not initialized")
	}

	return answerService.SearchAndAnswerStream(query, history, topK)
}

// SearchAndAnswerWithGraphStream 基于知识图谱的流式智能问答
func (c *AIClient) SearchAndAnswerWithGraphStream(query string, history []ChatMessage, topK int) (<-chan StreamChunk, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not available")
	}

	answerService := GetAnswerService()
	if answerService == nil {
		return nil, fmt.Errorf("answer service not initialized")
	}

	return answerService.SearchAndAnswerWithGraphStream(query, history, topK)
}

// BuildConversationContext 构建对话上下文字符串
func BuildConversationContext(history []ChatMessage) string {
	if len(history) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("对话历史：\n")
	for _, m := range history {
		if m.Role == "user" {
			sb.WriteString(fmt.Sprintf("用户：%s\n", m.Content))
		} else {
			sb.WriteString(fmt.Sprintf("助手：%s\n", m.Content))
		}
	}
	return sb.String()
}

// GetGraph 从数据库获取完整知识图谱数据
func (c *AIClient) GetGraph() (*AIGraphResponse, error) {
	// 从数据库获取知识点和关系
	points, rels, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil {
		return nil, fmt.Errorf("failed to get graph data: %w", err)
	}

	// 转换知识点为图节点
	nodes := make([]AIGraphNode, 0, len(points))
	for _, p := range points {
		nodes = append(nodes, AIGraphNode{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			DocumentID:  p.DocumentID,
		})
	}

	// 转换关系为图边
	edges := make([]AIGraphEdge, 0, len(rels))
	for _, r := range rels {
		edges = append(edges, AIGraphEdge{
			Source:       r.SourceID,
			Target:       r.TargetID,
			RelationType: r.RelationType,
			Description:  r.Description,
		})
	}

	return &AIGraphResponse{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// aiClient 全局 AI 客户端单例
var aiClient *AIClient

// InitAIClient 初始化全局 AI 客户端
func InitAIClient() {
	// 初始化向量服务
	InitVectorService()

	// 初始化抽取服务
	InitExtractionService()

	// 初始化问答服务
	InitAnswerService()

	aiClient = &AIClient{
		available: true,
	}

	log.Println("AI client initialized (local mode)")
}

// GetAIClient 获取 AI 客户端实例
func GetAIClient() *AIClient {
	return aiClient
}
