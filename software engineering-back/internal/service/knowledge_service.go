package service

import (
	"errors"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

// CreateKnowledgePoint 创建新的知识点
func CreateKnowledgePoint(req request.CreateKnowledgePointRequest) (uint, error) {
	kp := &entity.KnowledgePoint{
		Name:        req.Name,
		Description: req.Description,
		DocumentID:  req.DocumentID,
		Category:    req.Category,
	}
	if err := repository.CreateKnowledgePoint(kp); err != nil {
		return 0, err
	}
	return kp.ID, nil
}

// GetKnowledgePoint 根据 ID 获取知识点详情
func GetKnowledgePoint(id uint) (*response.KnowledgePointResponse, error) {
	kp, err := repository.FindKnowledgePointByID(id)
	if err != nil {
		return nil, errors.New("知识点不存在")
	}
	return &response.KnowledgePointResponse{
		ID:          kp.ID,
		Name:        kp.Name,
		Description: kp.Description,
		DocumentID:  kp.DocumentID,
		Category:    kp.Category,
		CreatedAt:   kp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   kp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateKnowledgePoint 更新知识点信息，仅更新非空字段
func UpdateKnowledgePoint(id uint, req request.UpdateKnowledgePointRequest) error {
	kp, err := repository.FindKnowledgePointByID(id)
	if err != nil {
		return errors.New("知识点不存在")
	}
	if req.Name != "" {
		kp.Name = req.Name
	}
	if req.Description != "" {
		kp.Description = req.Description
	}
	if req.Category != "" {
		kp.Category = req.Category
	}
	return repository.UpdateKnowledgePoint(kp)
}

// DeleteKnowledgePoint 删除知识点
func DeleteKnowledgePoint(id uint) error {
	_, err := repository.FindKnowledgePointByID(id)
	if err != nil {
		return errors.New("知识点不存在")
	}
	return repository.DeleteKnowledgePoint(id)
}

// ListKnowledgePoints 分页查询知识点列表，支持按关键词和文档 ID 过滤
func ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]response.KnowledgePointResponse, int64, error) {
	points, total, err := repository.ListKnowledgePoints(page, size, keyword, documentID)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.KnowledgePointResponse, len(points))
	for i, p := range points {
		list[i] = response.KnowledgePointResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			DocumentID:  p.DocumentID,
			Category:    p.Category,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

// --- 关系 ---

// CreateRelation 创建知识点之间的关系，会验证源和目标知识点是否存在
func CreateRelation(req request.CreateRelationRequest) (uint, error) {
	// 验证源和目标知识点存在
	if _, err := repository.FindKnowledgePointByID(req.SourceID); err != nil {
		return 0, errors.New("源知识点不存在")
	}
	if _, err := repository.FindKnowledgePointByID(req.TargetID); err != nil {
		return 0, errors.New("目标知识点不存在")
	}
	rel := &entity.KnowledgeRelation{
		SourceID:     req.SourceID,
		TargetID:     req.TargetID,
		RelationType: req.RelationType,
		Description:  req.Description,
	}
	if err := repository.CreateRelation(rel); err != nil {
		return 0, err
	}
	return rel.ID, nil
}

// UpdateRelation 更新知识点关系，仅更新非空字段
func UpdateRelation(id uint, req request.UpdateRelationRequest) error {
	rel, err := repository.FindRelationByID(id)
	if err != nil {
		return errors.New("关系不存在")
	}
	if req.RelationType != "" {
		rel.RelationType = req.RelationType
	}
	if req.Description != "" {
		rel.Description = req.Description
	}
	return repository.UpdateRelation(rel)
}

// DeleteRelation 删除知识点关系
func DeleteRelation(id uint) error {
	_, err := repository.FindRelationByID(id)
	if err != nil {
		return errors.New("关系不存在")
	}
	return repository.DeleteRelation(id)
}

// ListRelations 分页查询知识点关系列表，返回时附带源/目标知识点名称
func ListRelations(page, size int, pointID uint) ([]response.KnowledgeRelationResponse, int64, error) {
	rels, total, err := repository.ListRelations(page, size, pointID)
	if err != nil {
		return nil, 0, err
	}

	// 批量查询知识点名称
	pointIDs := make(map[uint]bool)
	for _, r := range rels {
		pointIDs[r.SourceID] = true
		pointIDs[r.TargetID] = true
	}
	names := make(map[uint]string)
	for id := range pointIDs {
		if kp, err := repository.FindKnowledgePointByID(id); err == nil {
			names[id] = kp.Name
		}
	}

	list := make([]response.KnowledgeRelationResponse, len(rels))
	for i, r := range rels {
		list[i] = response.KnowledgeRelationResponse{
			ID:           r.ID,
			SourceID:     r.SourceID,
			SourceName:   names[r.SourceID],
			TargetID:     r.TargetID,
			TargetName:   names[r.TargetID],
			RelationType: r.RelationType,
			Description:  r.Description,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

// ListKnowledgePointsSimple 简化的知识点列表查询（管理员使用）
func ListKnowledgePointsSimple(page, size int) ([]response.KnowledgePointResponse, int64, error) {
	return ListKnowledgePoints(page, size, "", 0)
}

// ListKnowledgeRelationsSimple 简化的关系列表查询（管理员使用）
func ListKnowledgeRelationsSimple(page, size int) ([]response.KnowledgeRelationResponse, int64, error) {
	return ListRelations(page, size, 0)
}
