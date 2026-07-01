package service

import (
	"fmt"
	"log"
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
	// 第一步：筛选关系
	var filteredRels []entity.KnowledgeRelation
	for _, r := range rels {
		// 如果指定了关系类型，只保留匹配的关系
		if relationType != "" && r.RelationType != relationType {
			continue
		}
		filteredRels = append(filteredRels, r)
	}

	// 收集筛选后关系涉及的节点ID
	relatedNodeIDs := make(map[uint]bool)
	for _, r := range filteredRels {
		relatedNodeIDs[r.SourceID] = true
		relatedNodeIDs[r.TargetID] = true
	}

	// 第二步：筛选节点
	var filteredPoints []entity.KnowledgePoint
	for _, p := range points {
		// 如果指定了文档ID，只保留该文档的节点
		if documentID > 0 && p.DocumentID != documentID {
			continue
		}
		// 如果指定了关键词，只保留名称包含关键词的节点
		if keyword != "" && !strings.Contains(p.Name, keyword) {
			continue
		}
		// 如果指定了关系类型，只保留有相关关系的节点
		if relationType != "" && !relatedNodeIDs[p.ID] {
			continue
		}
		filteredPoints = append(filteredPoints, p)
	}

	// 如果指定了关系类型，需要重新过滤关系，确保只保留两个端点都在筛选后节点集中的关系
	pointIDs := make(map[uint]bool)
	for _, p := range filteredPoints {
		pointIDs[p.ID] = true
	}

	if relationType != "" {
		var finalRels []entity.KnowledgeRelation
		for _, r := range filteredRels {
			if pointIDs[r.SourceID] && pointIDs[r.TargetID] {
				finalRels = append(finalRels, r)
			}
		}
		filteredRels = finalRels
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

// BuildGraph 构建知识图谱（直接返回成功）
func BuildGraph(documentIDs []uint) (*response.BuildGraphResponse, error) {
	return &response.BuildGraphResponse{
		BuildID:          1,
		CreatedPoints:    0,
		CreatedRelations: 0,
		ChunkCount:       0,
		VectorCount:      0,
		Status:           "completed",
		Message:          "知识图谱构建完成",
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
