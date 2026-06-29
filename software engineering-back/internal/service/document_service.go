package service

import (
	"io"
	"log"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
	apperrors "software_engineering/pkg/errors"
)

// UploadDocument 上传文档，提取文本内容保存到数据库（仅支持 .md 和 .txt）
func UploadDocument(userID uint, title, description string, filename string, fileSize int64, fileType string, contentReader io.Reader) (*response.DocumentResponse, error) {
	// 读取文件内容
	contentBytes, err := io.ReadAll(contentReader)
	if err != nil {
		return nil, err
	}

	// 读取文本内容（仅对文本文件）
	var fileContent string
	ext := strings.ToLower(fileType)
	if ext == ".md" || ext == ".txt" {
		fileContent = decodeContentWithCharset(contentBytes)
	}

	if title == "" {
		title = filename
	}

	doc := &entity.Document{
		UserID:      userID,
		Title:       title,
		Description: description,
		Filename:    filename,
		FilePath:    "", // 不再需要文件路径
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

// decodeContentWithCharset 检测字节流编码并转为 UTF-8 字符串。
// Windows 中文系统默认编码为 GBK，浏览器上传 .md/.txt 文件时会按原编码发送字节流。
// 此函数先检测是否为合法 UTF-8，若不是则尝试 GBK 解码，失败则回退为原始字节按 UTF-8 解析。
func decodeContentWithCharset(data []byte) string {
	if utf8.Valid(data) {
		return string(data)
	}
	// 尝试 GBK 解码
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(decoder, data)
	if err == nil {
		return string(utf8Bytes)
	}
	// 回退：强制按 UTF-8 解析，替换无效字节
	return strings.ToValidUTF8(string(data), "")
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

// buildKnowledgeGraph 为文档构建知识图谱
func buildKnowledgeGraph(documentID uint, content string) {
	extractionSvc := GetExtractionService()
	if extractionSvc == nil {
		log.Printf("Extraction service not initialized, skipping graph build for document %d", documentID)
		return
	}

	// 抽取知识点
	result, err := extractionSvc.ExtractKnowledgePoints(content, documentID)
	if err != nil {
		log.Printf("Failed to extract knowledge points for document %d: %v", documentID, err)
		return
	}
	log.Printf("Extracted %d points and %d relations for document %d", len(result.Points), len(result.Relations), documentID)

	// 构建图谱
	_, err = BuildGraph([]uint{documentID})
	if err != nil {
		log.Printf("Failed to build graph for document %d: %v", documentID, err)
		return
	}
	log.Printf("Built knowledge graph for document %d", documentID)
}

// GetDocument 获取文档详情，返回内容预览（前 200 字符）
func GetDocument(id uint) (*response.DocumentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeDocumentNotFound, "文档不存在")
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
		return nil, apperrors.New(apperrors.CodeDocumentNotFound, "文档不存在")
	}
	return &response.DocumentContentResponse{ID: doc.ID, Title: doc.Title, Content: doc.Content}, nil
}

// UpdateDocument 更新文档信息，仅更新非空字段
func UpdateDocument(id uint, req request.UpdateDocumentRequest) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return apperrors.New(apperrors.CodeDocumentNotFound, "文档不存在")
	}
	if req.Title != "" {
		doc.Title = req.Title
	}
	if req.Description != "" {
		doc.Description = req.Description
	}
	return repository.UpdateDocument(doc)
}

// DeleteDocument 删除文档，包含归属校验
func DeleteDocument(userID uint, id uint) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return apperrors.New(apperrors.CodeDocumentNotFound, "文档不存在")
	}
	// 归属校验：只能删除自己的文档
	if doc.UserID != userID {
		return apperrors.New(apperrors.CodeDocumentAccessDenied, "无权删除此文档")
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
		return apperrors.New(apperrors.CodeDocumentNotFound, "文档不存在")
	}

	doc.Status = req.Status
	doc.ReviewComment = req.Comment

	if err := repository.UpdateDocument(doc); err != nil {
		return err
	}

	// 审核通过时构建向量索引和知识图谱
	if req.Status == "approved" && doc.Content != "" {
		go buildVectorIndex(doc.ID, doc.Content)
		go buildKnowledgeGraph(doc.ID, doc.Content)
	}

	return nil
}

// DeleteDocumentAdmin 管理员删除文档（无需归属校验）
func DeleteDocumentAdmin(id uint) error {
	_, err := repository.FindDocumentByID(id)
	if err != nil {
		return apperrors.New(apperrors.CodeDocumentNotFound, "文档不存在")
	}
	return repository.DeleteDocument(id)
}
