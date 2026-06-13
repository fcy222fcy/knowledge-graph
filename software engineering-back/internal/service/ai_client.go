package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAIClient() *AIClient {
	return &AIClient{
		BaseURL: os.Getenv("AI_SERVICE_URL"),
		Client: &http.Client{
			Timeout: 30 * time.Second,
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
	DocumentID       uint   `json:"document_id"`
	DocumentTitle    string `json:"document_title"`
	CreatedPoints    int    `json:"created_points"`
	CreatedRelations int    `json:"created_relations"`
	ChunkCount       int    `json:"chunk_count"`
	VectorCount      int    `json:"vector_count"`
	Status           string `json:"status"`
	Message          string `json:"message"`
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
