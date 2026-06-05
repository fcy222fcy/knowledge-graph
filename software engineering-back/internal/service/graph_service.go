package service

import (
	"fmt"
	"strconv"
	"strings"

	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

func GetGraphData(documentID uint, keyword string, relationType string) (*response.GraphDataResponse, error) {
	// 优先从 Python AI 服务获取
	if aiClient.IsAvailable() {
		graphData, err := aiClient.GetGraph()
		if err == nil {
			return convertAIGraphToDTO(graphData, documentID, keyword, relationType), nil
		}
		// 降级到 Neo4j
	}

	// 从 Neo4j 获取
	points, rels, err := repository.GetAllGraphDataFromNeo4j()
	if err != nil {
		return nil, err
	}

	return filterAndConvertGraphData(points, rels, documentID, keyword, relationType), nil
}

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
		if !pointIDs[r.SourceID] && !pointIDs[r.TargetID] {
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

func convertAIGraphToDTO(graphData *AIGraphResponse, documentID uint, keyword string, relationType string) *response.GraphDataResponse {
	var nodes []response.GraphNode
	for _, n := range graphData.Nodes {
		if documentID > 0 && n.DocumentID != documentID {
			continue
		}
		if keyword != "" && !strings.Contains(n.Name, keyword) {
			continue
		}
		nodes = append(nodes, response.GraphNode{
			ID:          n.ID,
			Name:        n.Name,
			Description: n.Description,
			DocumentID:  n.DocumentID,
			Category:    n.Category,
		})
	}

	var edges []response.GraphEdge
	for _, e := range graphData.Edges {
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
			resp, err := aiClient.BuildGraph(AIBuildRequest{
				DocumentID: docID,
				Title:      doc.Title,
				Content:    doc.Content,
				Source:     "document",
			})
			if err == nil {
				totalPoints += resp.CreatedPoints
				totalRelations += resp.CreatedRelations
				totalChunks += resp.ChunkCount
				continue
			}
		}

		// 降级：使用简化的本地构建
		existingPoints, _ := repository.GetAllKnowledgePointsForGraph()
		var docPoints []entity.KnowledgePoint
		for _, p := range existingPoints {
			if p.DocumentID == docID {
				docPoints = append(docPoints, p)
			}
		}

		for i := 0; i < len(docPoints); i++ {
			for j := i + 1; j < len(docPoints); j++ {
				rel := &entity.KnowledgeRelation{
					SourceID:     docPoints[i].ID,
					TargetID:     docPoints[j].ID,
					RelationType: "RELATED",
					Description:  fmt.Sprintf("%s 与 %s 相关", docPoints[i].Name, docPoints[j].Name),
				}
				repository.CreateRelation(rel)
				totalRelations++
			}
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
