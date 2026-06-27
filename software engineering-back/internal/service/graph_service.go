package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

// GetGraphData 获取知识图谱数据，优先从 AI 服务获取，失败则降级到 Neo4j/MySQL
func GetGraphData(documentID uint, keyword string, relationType string) (*response.GraphDataResponse, error) {
	// 优先从 Python AI 服务获取
	if aiClient.IsAvailable() {
		graphData, err := aiClient.GetGraph()
		if err == nil {
			return convertAIGraphToDTO(graphData, documentID, keyword, relationType), nil
		}
		log.Printf("warning: AI graph service failed, degrading to Neo4j/MySQL: %v", err)
	} else {
		log.Println("info: AI graph service not available, using Neo4j/MySQL")
	}

	// 从 Neo4j 获取（如果不可用则自动降级到 MySQL）
	points, rels, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil {
		// Neo4j 失败时降级到 MySQL
		points, err = repository.GetAllKnowledgePointsForGraph()
		if err != nil {
			return nil, err
		}
		rels, err = repository.GetAllRelationsForGraph()
		if err != nil {
			return nil, err
		}
	}

	return filterAndConvertGraphData(points, rels, documentID, keyword, relationType), nil
}

// filterAndConvertGraphData 过滤节点和边，并转换为前端响应格式
func filterAndConvertGraphData(points []entity.KnowledgePoint, rels []entity.KnowledgeRelation, documentID uint, keyword string, relationType string) *response.GraphDataResponse {
	// 过滤
	var filteredPoints []entity.KnowledgePoint
	for _, p := range points {
		if documentID > 0 && p.DocumentID != documentID {
			continue
		}
		if keyword != "" && !strings.Contains(p.Name, keyword) {
			continue
		}
		filteredPoints = append(filteredPoints, p)
	}

	pointIDs := make(map[uint]bool)
	for _, p := range filteredPoints {
		pointIDs[p.ID] = true
	}

	var filteredRels []entity.KnowledgeRelation
	for _, r := range rels {
		// 边的 source 和 target 都必须在过滤后的节点集中
		if !pointIDs[r.SourceID] || !pointIDs[r.TargetID] {
			continue
		}
		if relationType != "" && r.RelationType != relationType {
			continue
		}
		filteredRels = append(filteredRels, r)
	}

	nodes := make([]response.GraphNode, len(filteredPoints))
	for i, p := range filteredPoints {
		nodes[i] = response.GraphNode{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			DocumentID:  p.DocumentID,
			Category:    p.Category,
		}
	}

	edges := make([]response.GraphEdge, len(filteredRels))
	for i, r := range filteredRels {
		edges[i] = response.GraphEdge{
			ID:           r.ID,
			Source:       r.SourceID,
			Target:       r.TargetID,
			RelationType: r.RelationType,
			Description:  r.Description,
		}
	}

	return &response.GraphDataResponse{
		Nodes: nodes,
		Edges: edges,
		Summary: response.GraphSummary{
			NodeCount: len(nodes),
			EdgeCount: len(edges),
		},
	}
}

// convertAIGraphToDTO 将 AI 服务返回的图数据转换为前端 DTO，支持按文档/关键词/关系类型过滤
func convertAIGraphToDTO(graphData *AIGraphResponse, documentID uint, keyword string, relationType string) *response.GraphDataResponse {
	nodes := make([]response.GraphNode, 0)
	filteredNodeIDs := make(map[uint]bool)
	for _, n := range graphData.Nodes {
		if documentID > 0 && n.DocumentID != documentID {
			continue
		}
		if keyword != "" && !strings.Contains(n.Name, keyword) {
			continue
		}
		filteredNodeIDs[n.ID] = true
		nodes = append(nodes, response.GraphNode{
			ID:          n.ID,
			Name:        n.Name,
			Description: n.Description,
			DocumentID:  n.DocumentID,
			Category:    n.Category,
		})
	}

	edges := make([]response.GraphEdge, 0)
	for _, e := range graphData.Edges {
		// 交叉验证：边的 source 和 target 都必须在过滤后的节点集中
		if !filteredNodeIDs[e.Source] || !filteredNodeIDs[e.Target] {
			continue
		}
		if relationType != "" && e.RelationType != relationType {
			continue
		}
		edges = append(edges, response.GraphEdge{
			Source:       e.Source,
			Target:       e.Target,
			RelationType: e.RelationType,
			Description:  e.Description,
		})
	}

	return &response.GraphDataResponse{
		Nodes: nodes,
		Edges: edges,
		Summary: response.GraphSummary{
			NodeCount: len(nodes),
			EdgeCount: len(edges),
		},
	}
}

// BuildGraph 根据文档 ID 列表构建知识图谱，优先调用 AI 服务，降级时使用本地链式关系构建
func BuildGraph(documentIDs []uint) (*response.BuildGraphResponse, error) {
	// 从 MySQL 读取文档内容
	var totalPoints, totalRelations, totalChunks int

	for _, docID := range documentIDs {
		doc, err := repository.FindDocumentByID(docID)
		if err != nil {
			continue
		}

		// 调用 Python AI 服务构建
		if aiClient.IsAvailable() {
			log.Printf("info: Calling AI service for document %d, content length: %d", docID, len(doc.Content))
			resp, err := aiClient.BuildGraph(AIBuildRequest{
				DocumentID: docID,
				Title:      doc.Title,
				Content:    doc.Content,
				Source:     "document",
			})
			if err == nil {
				log.Printf("info: AI build success for document %d: %d points, %d relations", docID, resp.CreatedPoints, resp.CreatedRelations)
				totalPoints += resp.CreatedPoints
				totalRelations += resp.CreatedRelations
				totalChunks += resp.ChunkCount

				// 优先写入 Neo4j（主存储），MySQL 作为备份
				// 建立 Python ID 到 MySQL ID 的映射
				idMapping := make(map[uint]uint)
				for _, p := range resp.Points {
					kp := &entity.KnowledgePoint{
						Name:        p.Name,
						Description: p.Description,
						DocumentID:  p.DocumentID,
						Category:    p.Category,
					}

					// 优先写入 Neo4j
					if err := repository.CreateKnowledgePointInNeo4j(kp); err != nil {
						log.Printf("warning: failed to save knowledge point to Neo4j (will use MySQL backup): %v", err)
					}

					// MySQL 作为备份存储
					if err := repository.CreateKnowledgePoint(kp); err != nil {
						log.Printf("warning: failed to save knowledge point to MySQL: %v", err)
					} else {
						idMapping[p.ID] = kp.ID
					}
				}
				for _, r := range resp.Relations {
					// 使用映射后的 ID
					sourceID := idMapping[r.Source]
					targetID := idMapping[r.Target]
					if sourceID == 0 || targetID == 0 {
						log.Printf("warning: skipping relation with unmapped IDs: source=%d, target=%d", r.Source, r.Target)
						continue
					}
					rel := &entity.KnowledgeRelation{
						SourceID:     sourceID,
						TargetID:     targetID,
						RelationType: r.RelationType,
						Description:  r.Description,
					}

					// 优先写入 Neo4j
					if err := repository.CreateRelationInNeo4j(rel); err != nil {
						log.Printf("warning: failed to save relation to Neo4j (will use MySQL backup): %v", err)
					}

					// MySQL 作为备份存储
					if err := repository.CreateRelation(rel); err != nil {
						log.Printf("warning: failed to save relation to MySQL: %v", err)
					}
				}
				continue
			}
			log.Printf("warning: AI build graph failed for document %d, degrading to local: %v", docID, err)
		} else {
			log.Printf("info: AI service not available, using local graph build for document %d", docID)
		}

		// 降级：使用简化的本地构建 - 只建立链式关系，不全连接
		existingPoints, _ := repository.GetAllKnowledgePointsForGraph()
		var docPoints []entity.KnowledgePoint
		for _, p := range existingPoints {
			if p.DocumentID == docID {
				docPoints = append(docPoints, p)
			}
		}

		// 只建立相邻节点的链式关系
		for i := 0; i < len(docPoints)-1; i++ {
			rel := &entity.KnowledgeRelation{
				SourceID:     docPoints[i].ID,
				TargetID:     docPoints[i+1].ID,
				RelationType: "RELATED",
				Description:  fmt.Sprintf("%s 与 %s 相关", docPoints[i].Name, docPoints[i+1].Name),
			}
			repository.CreateRelation(rel)
			totalRelations++
		}
		totalChunks += len(docPoints)
	}

	// 记录构建历史
	docIDsStr := make([]string, len(documentIDs))
	for i, id := range documentIDs {
		docIDsStr[i] = strconv.Itoa(int(id))
	}

	build := &entity.KnowledgeBuild{
		DocumentIDs:      strings.Join(docIDsStr, ","),
		CreatedPoints:    totalPoints,
		CreatedRelations: totalRelations,
		ChunkCount:       totalChunks,
		VectorCount:      totalChunks * 3,
		Status:           "completed",
		Message:          "知识图谱构建完成",
	}
	repository.CreateKnowledgeBuild(build)

	return &response.BuildGraphResponse{
		BuildID:          build.ID,
		CreatedPoints:    totalPoints,
		CreatedRelations: totalRelations,
		ChunkCount:       build.ChunkCount,
		VectorCount:      build.VectorCount,
		Status:           build.Status,
		Message:          build.Message,
	}, nil
}

// GetLatestBuildResult 获取最近一次知识图谱构建结果
func GetLatestBuildResult() (*response.BuildGraphResponse, error) {
	build, err := repository.GetLatestBuild()
	if err != nil {
		return nil, fmt.Errorf("暂无构建记录")
	}
	return &response.BuildGraphResponse{
		BuildID:          build.ID,
		CreatedPoints:    build.CreatedPoints,
		CreatedRelations: build.CreatedRelations,
		ChunkCount:       build.ChunkCount,
		VectorCount:      build.VectorCount,
		Status:           build.Status,
		Message:          build.Message,
	}, nil
}

// ListBuildHistory 分页获取知识图谱构建历史记录
func ListBuildHistory(page, size int) (*response.BuildHistoryResponse, error) {
	builds, total, err := repository.ListBuilds(page, size)
	if err != nil {
		return nil, err
	}
	list := make([]response.BuildGraphResponse, len(builds))
	for i, b := range builds {
		list[i] = response.BuildGraphResponse{
			BuildID:          b.ID,
			CreatedPoints:    b.CreatedPoints,
			CreatedRelations: b.CreatedRelations,
			ChunkCount:       b.ChunkCount,
			VectorCount:      b.VectorCount,
			Status:           b.Status,
			Message:          b.Message,
		}
	}
	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}
	return &response.BuildHistoryResponse{List: list, Total: total, Page: page, Size: size, TotalPage: totalPage}, nil
}
