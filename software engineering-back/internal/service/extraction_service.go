package service

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
	"software_engineering/pkg/ai"
	"software_engineering/pkg/config"
)

// ExtractionService 知识抽取服务
type ExtractionService struct {
	ollamaClient *ai.OllamaClient
}

// ExtractedPoint 抽取的知识点
type ExtractedPoint struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ExtractedRelation 抽取的关系
type ExtractedRelation struct {
	Source       string `json:"source"`
	Target       string `json:"target"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
}

// ExtractionResult 抽取结果
type ExtractionResult struct {
	Points    []ExtractedPoint    `json:"points"`
	Relations []ExtractedRelation `json:"relations"`
}

// extractionService 抽取服务单例
var extractionService *ExtractionService

// 软件工程本体 - 用于正则降级方案
var ONTOLOGY = map[string]struct {
	Keywords []string
}{
	"需求工程": {Keywords: []string{"需求", "需求分析", "需求获取", "需求规格", "需求验证", "用户需求", "系统需求"}},
	"设计": {Keywords: []string{"设计", "架构设计", "详细设计", "模块设计", "接口设计", "数据库设计"}},
	"编码": {Keywords: []string{"编码", "编程", "代码", "实现", "开发", "程序"}},
	"测试": {Keywords: []string{"测试", "单元测试", "集成测试", "系统测试", "验收测试", "测试用例"}},
	"维护": {Keywords: []string{"维护", "软件维护", "改正性维护", "适应性维护", "完善性维护"}},
	"项目管理": {Keywords: []string{"项目管理", "进度", "成本", "风险", "质量", "团队"}},
	"配置管理": {Keywords: []string{"配置管理", "版本控制", "变更管理", "基线"}},
	"质量保证": {Keywords: []string{"质量", "质量保证", "质量控制", "评审", "审计"}},
}

// 有效的分类
var validCategories = []string{
	"需求工程", "设计", "编码", "测试", "维护",
	"项目管理", "配置管理", "质量保证", "文档", "工具",
}

// 有效的关系类型
var validRelationTypes = []string{
	"RELATED", "DEPENDS_ON", "PART_OF", "IS_A", "EXAMPLE_OF", "USES", "IMPLEMENTS",
}

// InitExtractionService 初始化抽取服务
func InitExtractionService() {
	cfg := config.AppConfig
	ollamaClient := ai.NewOllamaClient(ai.OllamaConfig{
		BaseURL:        cfg.OllamaURL,
		Model:          cfg.OllamaModel,
		EmbeddingModel: cfg.OllamaEmbeddingModel,
	})

	extractionService = &ExtractionService{
		ollamaClient: ollamaClient,
	}

	log.Println("Extraction service initialized")
}

// GetExtractionService 获取抽取服务实例
func GetExtractionService() *ExtractionService {
	return extractionService
}

// ExtractKnowledgePoints 从文档内容中抽取知识点和关系
func (s *ExtractionService) ExtractKnowledgePoints(content string, documentID uint) (*ExtractionResult, error) {
	// 尝试使用 LLM 抽取
	result, err := s.extractWithLLM(content)
	if err != nil {
		log.Printf("LLM extraction failed, falling back to regex: %v", err)
		// 降级到正则抽取
		result = s.extractWithRegex(content)
	}

	// 验证和清理结果
	result = s.validateResult(result)

	return result, nil
}

// ExtractAndStoreToNeo4j 抽取知识点和关系后直接写入 Neo4j 和 MySQL
func (s *ExtractionService) ExtractAndStoreToNeo4j(content string, documentID uint) (*ExtractionResult, error) {
	// 1. 使用 LLM 抽取知识点和关系
	result, err := s.extractWithLLM(content)
	if err != nil {
		log.Printf("LLM extraction failed, falling back to regex: %v", err)
		result = s.extractWithRegex(content)
	}

	// 2. 验证和清理结果
	result = s.validateResult(result)

	// 3. 存储知识点到 MySQL + Neo4j
	for _, point := range result.Points {
		kp := &entity.KnowledgePoint{
			Name:        point.Name,
			Description: point.Description,
			DocumentID:  documentID,
			Category:    point.Category,
		}
		if err := repository.CreateKnowledgePoint(kp); err != nil {
			log.Printf("failed to store knowledge point '%s': %v", point.Name, err)
			continue
		}
	}

	// 4. 存储关系到 MySQL + Neo4j
	for _, rel := range result.Relations {
		// 查找源和目标知识点以获取 ID
		sourceKP, err := repository.FindKnowledgePointByName(rel.Source, documentID)
		if err != nil {
			log.Printf("failed to find source knowledge point '%s': %v", rel.Source, err)
			continue
		}
		targetKP, err := repository.FindKnowledgePointByName(rel.Target, documentID)
		if err != nil {
			log.Printf("failed to find target knowledge point '%s': %v", rel.Target, err)
			continue
		}
		kr := &entity.KnowledgeRelation{
			SourceID:     sourceKP.ID,
			TargetID:     targetKP.ID,
			Type:         rel.RelationType,
			RelationType: rel.RelationType,
			Description:  rel.Description,
		}
		if err := repository.CreateRelation(kr); err != nil {
			log.Printf("failed to store relation '%s' -> '%s': %v", rel.Source, rel.Target, err)
			continue
		}
	}

	return result, nil
}

// extractWithLLM 使用 LLM 抽取知识点
func (s *ExtractionService) extractWithLLM(content string) (*ExtractionResult, error) {
	// 限制内容长度
	if len(content) > 6000 {
		content = content[:6000]
	}

	systemPrompt := `你是一个软件工程知识图谱构建专家。请从给定的文档内容中提取知识点和它们之间的关系。

输出格式要求（严格 JSON）：
{
  "points": [
    {"name": "知识点名称", "description": "简短描述", "category": "分类"}
  ],
  "relations": [
    {"source": "源知识点", "target": "目标知识点", "relation_type": "关系类型", "description": "关系描述"}
  ]
}

分类必须是以下之一：需求工程、设计、编码、测试、维护、项目管理、配置管理、质量保证、文档、工具
关系类型必须是以下之一：RELATED、DEPENDS_ON、PART_OF、IS_A、EXAMPLE_OF、USES、IMPLEMENTS

请提取 10-30 个知识点和 5-20 个关系。只输出 JSON，不要其他内容。`

	userPrompt := fmt.Sprintf("请从以下文档内容中提取知识点和关系：\n\n%s", content)

	response, err := s.ollamaClient.Generate(userPrompt, systemPrompt, &ai.GenerateOptions{
		Temperature: 0.3,
		TopP:        0.9,
		NumPredict:  2048,
	})
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	var result ExtractionResult
	// 尝试提取 JSON 部分
	jsonStr := extractJSON(response)
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return &result, nil
}

// extractWithRegex 使用正则表达式降级抽取
func (s *ExtractionService) extractWithRegex(content string) *ExtractionResult {
	result := &ExtractionResult{
		Points:    make([]ExtractedPoint, 0),
		Relations: make([]ExtractedRelation, 0),
	}

	// 按段落分割内容
	paragraphs := strings.Split(content, "\n\n")

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if len(paragraph) < 10 {
			continue
		}

		// 检查是否匹配本体中的关键词
		for category, info := range ONTOLOGY {
			for _, keyword := range info.Keywords {
				if strings.Contains(paragraph, keyword) {
					// 提取包含关键词的句子作为知识点描述
					sentences := strings.Split(paragraph, "。")
					for _, sentence := range sentences {
						if strings.Contains(sentence, keyword) && len(sentence) > 10 {
							point := ExtractedPoint{
								Name:        keyword,
								Description: strings.TrimSpace(sentence),
								Category:    category,
							}
							// 避免重复
							if !pointExists(result.Points, point) {
								result.Points = append(result.Points, point)
							}
							break
						}
					}
					break
				}
			}
		}
	}

	// 生成简单的关联关系
	for i := 0; i < len(result.Points); i++ {
		for j := i + 1; j < len(result.Points); j++ {
			if result.Points[i].Category == result.Points[j].Category {
				relation := ExtractedRelation{
					Source:       result.Points[i].Name,
					Target:       result.Points[j].Name,
					RelationType: "RELATED",
					Description:  fmt.Sprintf("同属%s领域", result.Points[i].Category),
				}
				result.Relations = append(result.Relations, relation)
			}
		}
	}

	return result
}

// validateResult 验证和清理抽取结果
func (s *ExtractionService) validateResult(result *ExtractionResult) *ExtractionResult {
	// 验证知识点
	validPoints := make([]ExtractedPoint, 0)
	for _, point := range result.Points {
		if point.Name == "" || point.Description == "" {
			continue
		}
		// 验证分类
		if !isValidCategory(point.Category) {
			point.Category = "文档"
		}
		validPoints = append(validPoints, point)
	}
	result.Points = validPoints

	// 验证关系
	validRelations := make([]ExtractedRelation, 0)
	pointNames := make(map[string]bool)
	for _, p := range result.Points {
		pointNames[p.Name] = true
	}
	for _, rel := range result.Relations {
		if rel.Source == "" || rel.Target == "" {
			continue
		}
		// 检查源和目标是否存在
		if !pointNames[rel.Source] || !pointNames[rel.Target] {
			continue
		}
		// 验证关系类型
		if !isValidRelationType(rel.RelationType) {
			rel.RelationType = "RELATED"
		}
		validRelations = append(validRelations, rel)
	}
	result.Relations = validRelations

	return result
}

// pointExists 检查知识点是否已存在
func pointExists(points []ExtractedPoint, point ExtractedPoint) bool {
	for _, p := range points {
		if p.Name == point.Name && p.Category == point.Category {
			return true
		}
	}
	return false
}

// isValidCategory 检查分类是否有效
func isValidCategory(category string) bool {
	for _, c := range validCategories {
		if c == category {
			return true
		}
	}
	return false
}

// isValidRelationType 检查关系类型是否有效
func isValidRelationType(relType string) bool {
	for _, t := range validRelationTypes {
		if t == relType {
			return true
		}
	}
	return false
}

// extractJSON 从文本中提取 JSON 部分
func extractJSON(text string) string {
	// 尝试找到 JSON 块
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}

	// 如果没有找到，尝试清理文本
	re := regexp.MustCompile(`(?s)\{.*\}`)
	match := re.FindString(text)
	if match != "" {
		return match
	}

	return text
}

// GenerateChunkIndex 为文档生成分块索引
// 按 Markdown 标题结构分块，每个 chunk 是一个完整的章节
func (s *ExtractionService) GenerateChunkIndex(content string, documentID uint) []VectorMetadata {
	chunks := make([]VectorMetadata, 0)

	// 按行分割
	lines := strings.Split(content, "\n")

	// 标题正则：匹配 ##, ###, #### 等（二级及以下标题）
	headingRegex := regexp.MustCompile(`^(#{2,6})\s+(.+)$`)

	// 存储每个 chunk 的标题层级和内容
	type section struct {
		level   int
		heading string
		content strings.Builder
	}

	sections := make([]section, 0)
	currentSection := section{level: 0, heading: "文档开头"}

	for _, line := range lines {
		// 检查是否是标题行
		matches := headingRegex.FindStringSubmatch(line)
		if matches != nil {
			// 保存之前的 section
			if currentSection.content.Len() > 0 {
				sections = append(sections, currentSection)
			}
			// 开始新 section
			currentSection = section{
				level:   len(matches[1]),
				heading: matches[2],
			}
		}
		currentSection.content.WriteString(line)
		currentSection.content.WriteString("\n")
	}

	// 保存最后一个 section
	if currentSection.content.Len() > 0 {
		sections = append(sections, currentSection)
	}

	// 合并过小的 section，确保每个 chunk 有一定内容
	maxChunkSize := 1500 // 最大 chunk 大小

	mergedChunks := make([]string, 0)
	currentMerged := ""

	for _, sec := range sections {
		secContent := sec.content.String()
		secLen := len(secContent)

		// 如果当前 section 本身很大，直接保存
		if secLen > maxChunkSize {
			if currentMerged != "" {
				mergedChunks = append(mergedChunks, currentMerged)
				currentMerged = ""
			}
			mergedChunks = append(mergedChunks, secContent)
			continue
		}

		// 尝试合并
		if len(currentMerged)+secLen > maxChunkSize && currentMerged != "" {
			// 当前合并块已经达到上限，保存并开始新块
			mergedChunks = append(mergedChunks, currentMerged)
			currentMerged = secContent
		} else {
			// 继续合并
			if currentMerged != "" {
				currentMerged += "\n"
			}
			currentMerged += secContent
		}
	}

	// 保存最后一个合并块
	if currentMerged != "" {
		mergedChunks = append(mergedChunks, currentMerged)
	}

	// 转换为 VectorMetadata
	for _, chunk := range mergedChunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		chunks = append(chunks, VectorMetadata{
			ChunkText:  chunk,
			DocumentID: documentID,
		})
	}

	// 如果没有生成任何 chunk（纯文本无标题），使用原来的段落分割逻辑
	if len(chunks) == 0 {
		chunks = s.generateFallbackChunks(content, documentID)
	}

	return chunks
}

// generateFallbackChunks 无标题时的降级分块策略
func (s *ExtractionService) generateFallbackChunks(content string, documentID uint) []VectorMetadata {
	chunks := make([]VectorMetadata, 0)
	paragraphs := strings.Split(content, "\n\n")
	chunkSize := 500
	currentChunk := ""

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if len(currentChunk)+len(para) > chunkSize && currentChunk != "" {
			chunks = append(chunks, VectorMetadata{
				ChunkText:  currentChunk,
				DocumentID: documentID,
			})
			currentChunk = para
		} else {
			if currentChunk != "" {
				currentChunk += "\n\n"
			}
			currentChunk += para
		}
	}

	if currentChunk != "" {
		chunks = append(chunks, VectorMetadata{
			ChunkText:  currentChunk,
			DocumentID: documentID,
		})
	}

	return chunks
}

// Ensure extraction_service uses time package
var _ = time.Now
