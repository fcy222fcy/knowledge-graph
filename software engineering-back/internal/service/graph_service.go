package service

import (
	"fmt"
	"strconv"
	"strings"

	"software_engineering/internal/model/dto"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

func GetGraphData(documentID uint, keyword string, relationType string) (*dto.GraphDataResponse, error) {
	points, err := repository.GetAllKnowledgePointsForGraph()
	if err != nil {
		return nil, err
	}
	rels, err := repository.GetAllRelationsForGraph()
	if err != nil {
		return nil, err
	}

	// Filter points
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

	nodes := make([]dto.GraphNode, len(filteredPoints))
	for i, p := range filteredPoints {
		nodes[i] = dto.GraphNode{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			DocumentID:  p.DocumentID,
			Category:    p.Category,
		}
	}

	edges := make([]dto.GraphEdge, len(filteredRels))
	for i, r := range filteredRels {
		edges[i] = dto.GraphEdge{
			ID:           r.ID,
			Source:       r.SourceID,
			Target:       r.TargetID,
			RelationType: r.RelationType,
			Description:  r.Description,
		}
	}

	return &dto.GraphDataResponse{
		Nodes: nodes,
		Edges: edges,
		Summary: dto.GraphSummary{
			NodeCount: len(nodes),
			EdgeCount: len(edges),
		},
	}, nil
}

func BuildGraph(documentIDs []uint) (*dto.BuildGraphResponse, error) {
	// Get all knowledge points
	existingPoints, _ := repository.GetAllKnowledgePointsForGraph()

	createdPoints := 0
	createdRelations := 0

	// Find points related to the given documents
	var docPoints []entity.KnowledgePoint
	for _, p := range existingPoints {
		for _, docID := range documentIDs {
			if p.DocumentID == docID {
				docPoints = append(docPoints, p)
			}
		}
	}

	// Create relations between points from the same document
	for i := 0; i < len(docPoints); i++ {
		for j := i + 1; j < len(docPoints); j++ {
			rel := &entity.KnowledgeRelation{
				SourceID:     docPoints[i].ID,
				TargetID:     docPoints[j].ID,
				RelationType: "RELATED",
				Description:  fmt.Sprintf("%s 与 %s 相关", docPoints[i].Name, docPoints[j].Name),
			}
			repository.CreateRelation(rel)
			createdRelations++
		}
	}

	docIDsStr := make([]string, len(documentIDs))
	for i, id := range documentIDs {
		docIDsStr[i] = strconv.Itoa(int(id))
	}

	build := &entity.KnowledgeBuild{
		DocumentIDs:      strings.Join(docIDsStr, ","),
		CreatedPoints:    createdPoints,
		CreatedRelations: createdRelations,
		ChunkCount:       len(docPoints),
		VectorCount:      len(docPoints) * 3,
		Status:           "completed",
		Message:          "知识图谱构建完成",
	}
	repository.CreateKnowledgeBuild(build)

	return &dto.BuildGraphResponse{
		BuildID:          build.ID,
		CreatedPoints:    createdPoints,
		CreatedRelations: createdRelations,
		ChunkCount:       build.ChunkCount,
		VectorCount:      build.VectorCount,
		Status:           build.Status,
		Message:          build.Message,
	}, nil
}

func GetLatestBuildResult() (*dto.BuildGraphResponse, error) {
	build, err := repository.GetLatestBuild()
	if err != nil {
		return nil, fmt.Errorf("暂无构建记录")
	}
	return &dto.BuildGraphResponse{
		BuildID:          build.ID,
		CreatedPoints:    build.CreatedPoints,
		CreatedRelations: build.CreatedRelations,
		ChunkCount:       build.ChunkCount,
		VectorCount:      build.VectorCount,
		Status:           build.Status,
		Message:          build.Message,
	}, nil
}

func ListBuildHistory(page, size int) (*dto.BuildHistoryResponse, error) {
	builds, total, err := repository.ListBuilds(page, size)
	if err != nil {
		return nil, err
	}
	list := make([]dto.BuildGraphResponse, len(builds))
	for i, b := range builds {
		list[i] = dto.BuildGraphResponse{
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
	return &dto.BuildHistoryResponse{List: list, Total: total, Page: page, Size: size, TotalPage: totalPage}, nil
}
