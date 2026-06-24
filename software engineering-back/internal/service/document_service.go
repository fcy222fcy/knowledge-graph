package service

import (
	"context"
	"errors"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
	"software_engineering/pkg/config"
	"software_engineering/pkg/storage"
)

// minioClient MinIO 客户端单例
var minioClient *storage.MinIOClient

// InitMinIOClient 初始化 MinIO 客户端
func InitMinIOClient() error {
	cfg := config.AppConfig
	minioCfg := storage.MinIOConfig{
		Endpoint:  cfg.MinIOEndpoint,
		AccessKey: cfg.MinIOAccessKey,
		SecretKey: cfg.MinIOSecretKey,
		Bucket:    cfg.MinIOBucket,
		UseSSL:    cfg.MinIOUseSSL,
	}
	var err error
	minioClient, err = storage.NewMinIOClient(minioCfg)
	return err
}

// GetMinIOClient 获取 MinIO 客户端
func GetMinIOClient() *storage.MinIOClient {
	return minioClient
}

// UploadDocument 上传文档到 MinIO，保存文件并提取文本内容（仅支持 .md 和 .txt）
func UploadDocument(userID uint, title, description string, filename string, fileSize int64, fileType string, contentReader io.Reader) (*response.DocumentResponse, error) {
	if minioClient == nil {
		return nil, errors.New("MinIO client not initialized")
	}

	// 安全处理文件名（清理路径穿越字符）
	safeName := filepath.Base(filename)
	ctx := context.Background()

	// 生成唯一的 object name（用户ID/日期/文件名）
	objectName := filepath.Join(
		config.AppConfig.MinIOBucket,
		safeName,
	)

	// 读取文件内容到缓冲区（用于后续提取文本）
	contentBytes, err := io.ReadAll(contentReader)
	if err != nil {
		return nil, err
	}

	// 上传到 MinIO
	err = minioClient.UploadFile(ctx, objectName, strings.NewReader(string(contentBytes)), int64(len(contentBytes)), fileType)
	if err != nil {
		return nil, err
	}

	// 读取文本内容（仅对文本文件）
	var fileContent string
	ext := strings.ToLower(fileType)
	if ext == ".md" || ext == ".txt" {
		fileContent = string(contentBytes)
	}

	if title == "" {
		title = filename
	}

	doc := &entity.Document{
		UserID:      userID,
		Title:       title,
		Description: description,
		Filename:    filename,
		FilePath:    objectName,
		FileSize:    fileSize,
		FileType:    fileType,
		Content:     fileContent,
		Status:      "pending",
	}
	if err := repository.CreateDocument(doc); err != nil {
		return nil, err
	}

	// 构建向量索引（仅对文本文件）
	if fileContent != "" {
		go buildVectorIndex(doc.ID, fileContent)
	}

	return &response.DocumentResponse{
		ID:          doc.ID,
		Title:       doc.Title,
		Description: doc.Description,
		Filename:    doc.Filename,
		FileSize:    doc.FileSize,
		FileType:    doc.FileType,
		Status:      doc.Status,
		CreatedAt:   doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// buildVectorIndex 为文档构建向量索引
func buildVectorIndex(documentID uint, content string) {
	extractionSvc := &ExtractionService{}
	chunks := extractionSvc.GenerateChunkIndex(content, documentID)

	vecSvc := GetVectorService()
	if vecSvc == nil {
		log.Printf("Vector service not initialized, skipping index build for document %d", documentID)
		return
	}

	for i, chunk := range chunks {
		embedding, err := vecSvc.Embed(chunk.ChunkText)
		if err != nil {
			log.Printf("Failed to embed chunk %d for document %d: %v", i, documentID, err)
			continue
		}
		if err := vecSvc.AddVector(embedding, chunk); err != nil {
			log.Printf("Failed to add vector for chunk %d: %v", i, err)
			continue
		}
	}

	if err := vecSvc.SaveIndex(); err != nil {
		log.Printf("Failed to save vector index: %v", err)
	} else {
		log.Printf("Built vector index for document %d: %d chunks", documentID, len(chunks))
	}
}

// GetDocumentDownloadURL 获取文档的下载 URL
func GetDocumentDownloadURL(id uint) (string, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return "", errors.New("文档不存在")
	}

	if minioClient == nil {
		return "", errors.New("MinIO client not initialized")
	}

	ctx := context.Background()
	presignedURL, err := minioClient.GetPresignedURL(ctx, doc.FilePath, 2*time.Hour)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// GetDocument 获取文档详情，返回内容预览（前 200 字符）
func GetDocument(id uint) (*response.DocumentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	preview := doc.Content
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	return &response.DocumentResponse{
		ID:             doc.ID,
		Title:          doc.Title,
		Description:    doc.Description,
		Filename:       doc.Filename,
		FileSize:       doc.FileSize,
		FileType:       doc.FileType,
		Status:         doc.Status,
		ContentPreview: preview,
		CreatedAt:      doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// GetDocumentContent 获取文档完整内容
func GetDocumentContent(id uint) (*response.DocumentContentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	return &response.DocumentContentResponse{ID: doc.ID, Title: doc.Title, Content: doc.Content}, nil
}

// UpdateDocument 更新文档信息，仅更新非空字段
func UpdateDocument(id uint, req request.UpdateDocumentRequest) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	if req.Title != "" {
		doc.Title = req.Title
	}
	if req.Description != "" {
		doc.Description = req.Description
	}
	return repository.UpdateDocument(doc)
}

// DeleteDocument 删除文档，包含归属校验和 MinIO 文件清理
func DeleteDocument(userID uint, id uint) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	// 归属校验：只能删除自己的文档
	if doc.UserID != userID {
		return errors.New("无权删除此文档")
	}
	// 删除 MinIO 文件
	if doc.FilePath != "" && minioClient != nil {
		ctx := context.Background()
		_ = minioClient.DeleteFile(ctx, doc.FilePath)
	}
	return repository.DeleteDocument(id)
}

// ListDocuments 分页获取所有文档列表，支持按关键词和状态过滤
func ListDocuments(page, size int, keyword, status string) ([]response.DocumentResponse, int64, error) {
	docs, total, err := repository.ListDocuments(page, size, keyword, status)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.DocumentResponse, len(docs))
	for i, d := range docs {
		list[i] = response.DocumentResponse{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			Filename:    d.Filename,
			FileSize:    d.FileSize,
			FileType:    d.FileType,
			Status:      d.Status,
			CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

// ListUserDocuments 分页获取指定用户的文档列表
func ListUserDocuments(userID uint, page, size int, keyword, status string) ([]response.DocumentResponse, int64, error) {
	docs, total, err := repository.ListDocumentsByUser(userID, page, size, keyword, status)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.DocumentResponse, len(docs))
	for i, d := range docs {
		list[i] = response.DocumentResponse{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			Filename:    d.Filename,
			FileSize:    d.FileSize,
			FileType:    d.FileType,
			Status:      d.Status,
			CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

// ListDocumentsAdmin 管理员获取所有文档列表
func ListDocumentsAdmin(page, size int, keyword string) ([]response.DocumentResponse, int64, error) {
	docs, total, err := repository.ListDocumentsAdmin(page, size, keyword)
	if err != nil {
		return nil, 0, err
	}
	list := make([]response.DocumentResponse, len(docs))
	for i, d := range docs {
		list[i] = response.DocumentResponse{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			Filename:    d.Filename,
			FileSize:    d.FileSize,
			FileType:    d.FileType,
			Status:      d.Status,
			CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

// ReviewDocument 教师/管理员审核文档，更新审核状态和意见
func ReviewDocument(id uint, req request.ReviewDocumentRequest) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}

	doc.Status = req.Status
	doc.ReviewComment = req.Comment

	if err := repository.UpdateDocument(doc); err != nil {
		return err
	}

	// 审核通过时构建向量索引
	if req.Status == "approved" && doc.Content != "" {
		go buildVectorIndex(doc.ID, doc.Content)
	}

	return nil
}

// DeleteDocumentAdmin 管理员删除文档（无需归属校验）
func DeleteDocumentAdmin(id uint) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	// 删除 MinIO 文件
	if doc.FilePath != "" && minioClient != nil {
		ctx := context.Background()
		_ = minioClient.DeleteFile(ctx, doc.FilePath)
	}
	return repository.DeleteDocument(id)
}
