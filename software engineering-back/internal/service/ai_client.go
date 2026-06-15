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

type AIClient struct {
	BaseURL string
	Client  *http.Client
}

type AIAnswerRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

type AIAnswerSource struct {
	DocumentID    int    `json:"document_id"`
	DocumentTitle string `json:"document_title"`
	Content       string `json:"content"`
}

type AIAnswerResponse struct {
	Answer     string           `json:"answer"`
	Confidence float64          `json:"confidence"`
	Sources    []AIAnswerSource `json:"sources"`
}

// ChatMessage 对话消息，用于传递历史上下文
type ChatMessage struct {
	Role    string `json:"role"`    // user / assistant
	Content string `json:"content"`
}

// AIAnswerWithHistoryRequest 带历史上下文的请求
type AIAnswerWithHistoryRequest struct {
	Query   string        `json:"query"`
	History []ChatMessage `json:"history"`
	TopK    int           `json:"top_k"`
}

func NewAIClient() *AIClient {
	return &AIClient{
		BaseURL: os.Getenv("AI_SERVICE_URL"),
		Client: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

type AIBuildRequest struct {
	DocumentID uint   `json:"document_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Source     string `json:"source"`
}

type AIBuildResponse struct {
	DocumentID       uint          `json:"document_id"`
	DocumentTitle    string        `json:"document_title"`
	CreatedPoints    int           `json:"created_points"`
	CreatedRelations int           `json:"created_relations"`
	ChunkCount       int           `json:"chunk_count"`
	VectorCount      int           `json:"vector_count"`
	Status           string        `json:"status"`
	Message          string        `json:"message"`
	Points           []AIGraphNode `json:"points"`
	Relations        []AIGraphEdge `json:"relations"`
}

type AISearchRequest struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

type AISearchResult struct {
	ChunkText         string  `json:"chunk_text"`
	Score             float64 `json:"score"`
	DocumentID        int     `json:"document_id"`
	KnowledgePointIDs []int   `json:"knowledge_point_ids"`
}

type AISearchResponse struct {
	Results []AISearchResult `json:"results"`
}

type AIGraphResponse struct {
	Nodes []AIGraphNode `json:"nodes"`
	Edges []AIGraphEdge `json:"edges"`
}

type AIGraphNode struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	DocumentID  uint   `json:"document_id"`
}

type AIGraphEdge struct {
	Source       uint   `json:"source"`
	Target       uint   `json:"target"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
}

func (c *AIClient) IsAvailable() bool {
	return c.BaseURL != ""
}

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

var aiClient *AIClient

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