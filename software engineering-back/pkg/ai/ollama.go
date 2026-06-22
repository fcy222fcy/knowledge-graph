package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient Ollama API 客户端
type OllamaClient struct {
	BaseURL        string
	Model          string
	EmbeddingModel string
	Client         *http.Client
}

// OllamaConfig Ollama 配置
type OllamaConfig struct {
	BaseURL        string
	Model          string
	EmbeddingModel string
}

// GenerateRequest 文本生成请求
type GenerateRequest struct {
	Model   string          `json:"model"`
	Prompt  string          `json:"prompt"`
	System  string          `json:"system,omitempty"`
	Stream  bool            `json:"stream"`
	Think   bool            `json:"think"`
	Options *GenerateOptions `json:"options,omitempty"`
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// GenerateResponse 文本生成响应
type GenerateResponse struct {
	Response string `json:"response"`
	Thinking string `json:"thinking,omitempty"`
}

// GenerateStreamResponse 流式生成响应
type GenerateStreamResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Thinking string `json:"thinking,omitempty"`
}

// EmbeddingRequest 嵌入向量请求
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse 嵌入向量响应
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// NewOllamaClient 创建 Ollama 客户端
func NewOllamaClient(cfg OllamaConfig) *OllamaClient {
	return &OllamaClient{
		BaseURL:        cfg.BaseURL,
		Model:          cfg.Model,
		EmbeddingModel: cfg.EmbeddingModel,
		Client: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

// Generate 调用 Ollama 生成文本
func (c *OllamaClient) Generate(prompt, system string, options *GenerateOptions) (string, error) {
	req := GenerateRequest{
		Model:   c.Model,
		Prompt:  prompt,
		System:  system,
		Stream:  false,
		Think:   false,
		Options: options,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.Client.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to call ollama generate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama generate failed: %s", string(respBody))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 优先返回 response，如果为空则返回 thinking
	if result.Response != "" {
		return result.Response, nil
	}
	return result.Thinking, nil
}

// GenerateStream 流式生成文本，返回一个 channel 用于接收生成的 token
func (c *OllamaClient) GenerateStream(prompt, system string, options *GenerateOptions) (<-chan GenerateStreamResponse, error) {
	req := GenerateRequest{
		Model:   c.Model,
		Prompt:  prompt,
		System:  system,
		Stream:  true,
		Think:   false,
		Options: options,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 使用不带超时的 client 进行长连接
	httpClient := &http.Client{
		Timeout: 0, // 无超时，流式连接
	}

	resp, err := httpClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama generate stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("ollama generate stream failed: %s", string(respBody))
	}

	// 创建 channel 用于传输流式数据
	ch := make(chan GenerateStreamResponse, 100)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var result GenerateStreamResponse
			if err := json.Unmarshal([]byte(line), &result); err != nil {
				continue
			}

			ch <- result

			if result.Done {
				break
			}
		}
	}()

	return ch, nil
}

// Embed 获取文本的嵌入向量
func (c *OllamaClient) Embed(text string) ([]float64, error) {
	req := EmbeddingRequest{
		Model:  c.EmbeddingModel,
		Prompt: text,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.Client.Post(c.BaseURL+"/api/embeddings", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama embeddings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama embeddings failed: %s", string(respBody))
	}

	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Embedding, nil
}

// IsAvailable 检查 Ollama 服务是否可用
func (c *OllamaClient) IsAvailable() bool {
	resp, err := c.Client.Get(c.BaseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
