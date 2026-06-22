package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AIClient Python AI 服务客户端，负责调用知识图谱构建、语义搜索和智能问答接口
type AIClient struct {
	BaseURL string      // AI 服务基础URL
	Client  *http.Client // HTTP 客户端
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

// NewAIClient 创建 AI 客户端，从 AI_SERVICE_URL 环境变量读取服务地址
func NewAIClient() *AIClient {
	return &AIClient{
		BaseURL: os.Getenv("AI_SERVICE_URL"),
		Client: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
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
	return c.BaseURL != ""
}

// BuildGraph 调用 AI 服务构建知识图谱，返回创建的节点和边数量
func (c *AIClient) BuildGraph(req AIBuildRequest) (*AIBuildResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not configured")
	}

	body, _ := json.Marshal(req)
	resp, err := c.Client.Post(c.BaseURL+"/build", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call AI build: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI build failed: %s", string(body))
	}

	var result AIBuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI build response: %w", err)
	}
	return &result, nil
}

// Search 调用 AI 服务进行语义搜索，返回与查询最相关的文档片段
func (c *AIClient) Search(query string, topK int) (*AISearchResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not configured")
	}

	req := AISearchRequest{Query: query, TopK: topK}
	body, _ := json.Marshal(req)
	resp, err := c.Client.Post(c.BaseURL+"/search", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call AI search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI search failed: %s", string(body))
	}

	var result AISearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI search response: %w", err)
	}
	return &result, nil
}

// SearchAndAnswer 调用 AI 服务进行语义搜索并生成回答（不带历史上下文）
func (c *AIClient) SearchAndAnswer(query string, topK int) (*AIAnswerResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not configured")
	}

	req := AIAnswerRequest{Query: query, TopK: topK}
	body, _ := json.Marshal(req)
	resp, err := c.Client.Post(c.BaseURL+"/search_and_answer", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call AI search_and_answer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI search_and_answer failed: %s", string(body))
	}

	var result AIAnswerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI search_and_answer response: %w", err)
	}
	return &result, nil
}

// SearchAndAnswerWithHistory 带对话历史的智能问答
func (c *AIClient) SearchAndAnswerWithHistory(query string, history []ChatMessage, topK int) (*AIAnswerResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not configured")
	}

	req := AIAnswerWithHistoryRequest{Query: query, History: history, TopK: topK}
	body, _ := json.Marshal(req)
	resp, err := c.Client.Post(c.BaseURL+"/search_and_answer", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call AI search_and_answer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI search_and_answer failed: %s", string(body))
	}

	var result AIAnswerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI search_and_answer response: %w", err)
	}
	return &result, nil
}

// SearchAndAnswerWithGraph 基于知识图谱的智能问答
func (c *AIClient) SearchAndAnswerWithGraph(query string, history []ChatMessage, topK int) (*AIAnswerWithGraphResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not configured")
	}

	req := AIAnswerWithHistoryRequest{Query: query, History: history, TopK: topK}
	body, _ := json.Marshal(req)

	resp, err := c.Client.Post(c.BaseURL+"/search_and_answer_with_graph", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call AI search_and_answer_with_graph: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI search_and_answer_with_graph failed: %s", string(body))
	}

	var result AIAnswerWithGraphResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI search_and_answer_with_graph response: %w", err)
	}
	return &result, nil
}

// BuildConversationContext 构建对话上下文字符串，用于本地降级时拼接 prompt
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

// GetGraph 从 AI 服务获取完整知识图谱数据（节点和边）
func (c *AIClient) GetGraph() (*AIGraphResponse, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("AI service not configured")
	}

	resp, err := c.Client.Get(c.BaseURL + "/graph")
	if err != nil {
		return nil, fmt.Errorf("failed to call AI graph: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI graph failed: %s", string(body))
	}

	var result AIGraphResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI graph response: %w", err)
	}
	return &result, nil
}

// aiClient 全局 AI 客户端单例
var aiClient *AIClient

// InitAIClient 初始化全局 AI 客户端，日志输出连接状态
func InitAIClient() {
	aiClient = NewAIClient()
	if aiClient.IsAvailable() {
		log.Printf("AI client initialized with URL: %s", aiClient.BaseURL)
	} else {
		log.Println("AI client: AI_SERVICE_URL not configured, AI features disabled")
	}
}

// GetAIClient 获取 AI 客户端实例
func GetAIClient() *AIClient {
	return aiClient
}