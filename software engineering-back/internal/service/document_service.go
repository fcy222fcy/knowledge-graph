package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

const uploadDir = "./uploads"

func UploadDocument(title, description string, filename string, fileSize int64, fileType string, contentReader io.Reader) (*dto.DocumentResponse, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	content, err := io.Copy(out, contentReader)
	if err != nil {
		return nil, err
	}

	// 读取文本内容（仅对文本文件）
	var fileContent string
	ext := strings.ToLower(fileType)
	if ext == ".md" || ext == ".txt" {
		data, err := os.ReadFile(filePath)
		if err == nil {
			fileContent = string(data)
		}
	}

	if title == "" {
		title = filename
	}

	doc := &model.Document{
		Title:       title,
		Description: description,
		Filename:    filename,
		FilePath:    filePath,
		FileSize:    content,
		FileType:    fileType,
		Content:     fileContent,
		Status:      "completed",
	}
	if err := repository.CreateDocument(doc); err != nil {
		return nil, err
	}

	return &dto.DocumentResponse{
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

func GetDocument(id uint) (*dto.DocumentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	preview := doc.Content
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	return &dto.DocumentResponse{
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

func GetDocumentContent(id uint) (*dto.DocumentContentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	return &dto.DocumentContentResponse{ID: doc.ID, Title: doc.Title, Content: doc.Content}, nil
}

func UpdateDocument(id uint, req dto.UpdateDocumentRequest) error {
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

func DeleteDocument(id uint) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	// 删除物理文件
	if doc.FilePath != "" {
		os.Remove(doc.FilePath)
	}
	return repository.DeleteDocument(id)
}

func ListDocuments(page, size int, keyword, status string) ([]dto.DocumentResponse, int64, error) {
	docs, total, err := repository.ListDocuments(page, size, keyword, status)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.DocumentResponse, len(docs))
	for i, d := range docs {
		list[i] = dto.DocumentResponse{
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
