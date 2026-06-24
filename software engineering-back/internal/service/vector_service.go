package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sync"

	"software_engineering/pkg/ai"
	"software_engineering/pkg/config"
)

// VectorService 向量索引服务
type VectorService struct {
	ollamaClient *ai.OllamaClient
	vectors      [][]float64
	metadata     []VectorMetadata
	dimension    int
	mu           sync.RWMutex
	indexFile    string
}

// VectorMetadata 向量元数据
type VectorMetadata struct {
	ChunkText      string `json:"chunk_text"`
	DocumentID     uint   `json:"document_id"`
	KnowledgePointIDs []uint `json:"knowledge_point_ids"`
}

// VectorIndex 向量索引持久化结构
type VectorIndex struct {
	Vectors   [][]float64       `json:"vectors"`
	Metadata  []VectorMetadata  `json:"metadata"`
	Dimension int               `json:"dimension"`
}

// vectorService 向量服务单例
var vectorService *VectorService

// InitVectorService 初始化向量服务
func InitVectorService() {
	cfg := config.AppConfig
	ollamaClient := ai.NewOllamaClient(ai.OllamaConfig{
		BaseURL:        cfg.OllamaURL,
		Model:          cfg.OllamaModel,
		EmbeddingModel: cfg.OllamaEmbeddingModel,
	})

	vectorService = &VectorService{
		ollamaClient: ollamaClient,
		vectors:      make([][]float64, 0),
		metadata:     make([]VectorMetadata, 0),
		dimension:    768,
		indexFile:    "data/vector_index.json",
	}

	// 尝试加载已有索引
	if err := vectorService.loadIndex(); err != nil {
		log.Printf("Warning: failed to load vector index: %v", err)
	}

	log.Println("Vector service initialized")
}

// GetVectorService 获取向量服务实例
func GetVectorService() *VectorService {
	return vectorService
}

// Embed 获取文本的嵌入向量
func (s *VectorService) Embed(text string) ([]float64, error) {
	embedding, err := s.ollamaClient.Embed(text)
	if err != nil {
		return nil, err
	}

	// L2 归一化
	norm := 0.0
	for _, v := range embedding {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range embedding {
			embedding[i] /= norm
		}
	}

	return embedding, nil
}

// AddVector 添加向量到索引
func (s *VectorService) AddVector(vector []float64, meta VectorMetadata) error {
	if len(vector) != s.dimension {
		return fmt.Errorf("vector dimension mismatch: expected %d, got %d", s.dimension, len(vector))
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors = append(s.vectors, vector)
	s.metadata = append(s.metadata, meta)

	return nil
}

// Search 搜索最相似的向量
func (s *VectorService) Search(query string, topK int) ([]SearchResult, error) {
	queryVector, err := s.Embed(query)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.vectors) == 0 {
		return []SearchResult{}, nil
	}

	// 计算余弦相似度（向量已归一化，点积即余弦相似度）
	scores := make([]float64, len(s.vectors))
	for i, vec := range s.vectors {
		score := 0.0
		for j, v := range vec {
			score += v * queryVector[j]
		}
		scores[i] = score
	}

	// 获取 topK 结果
	type indexScore struct {
		index int
		score float64
	}
	results := make([]indexScore, 0, topK)
	for i, score := range scores {
		if len(results) < topK {
			results = append(results, indexScore{i, score})
			// 简单插入排序
			for j := len(results) - 1; j > 0; j-- {
				if results[j].score > results[j-1].score {
					results[j], results[j-1] = results[j-1], results[j]
				} else {
					break
				}
			}
		} else if score > results[topK-1].score {
			results[topK-1] = indexScore{i, score}
			for j := len(results) - 1; j > 0; j-- {
				if results[j].score > results[j-1].score {
					results[j], results[j-1] = results[j-1], results[j]
				} else {
					break
				}
			}
		}
	}

	// 转换结果并过滤低相关性
	minScore := 0.5 // 最低相似度阈值
	searchResults := make([]SearchResult, 0, len(results))
	for _, r := range results {
		if r.score >= minScore {
			searchResults = append(searchResults, SearchResult{
				Metadata: s.metadata[r.index],
				Score:    r.score,
			})
		}
	}

	return searchResults, nil
}

// SearchResult 搜索结果
type SearchResult struct {
	Metadata VectorMetadata
	Score    float64
}

// SaveIndex 保存索引到文件
func (s *VectorService) SaveIndex() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	index := VectorIndex{
		Vectors:   s.vectors,
		Metadata:  s.metadata,
		Dimension: s.dimension,
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.indexFile, data, 0644)
}

// loadIndex 从文件加载索引
func (s *VectorService) loadIndex() error {
	data, err := os.ReadFile(s.indexFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var index VectorIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return err
	}

	// 验证维度
	if index.Dimension != s.dimension {
		log.Printf("Vector index dimension mismatch: expected %d, got %d, clearing index", s.dimension, index.Dimension)
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors = index.Vectors
	s.metadata = index.Metadata
	log.Printf("Loaded vector index: %d vectors", len(s.vectors))
	return nil
}

// ClearIndex 清空索引
func (s *VectorService) ClearIndex() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors = make([][]float64, 0)
	s.metadata = make([]VectorMetadata, 0)
}

// GetSize 获取索引大小
func (s *VectorService) GetSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.vectors)
}
