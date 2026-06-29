package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CreateDocument 创建文档记录
func CreateDocument(doc *entity.Document) error {
	return database.DB.Create(doc).Error
}

// FindDocumentByID 根据 ID 查找文档
func FindDocumentByID(id uint) (*entity.Document, error) {
	var doc entity.Document
	err := database.DB.First(&doc, id).Error
	return &doc, err
}

// UpdateDocument 更新文档信息
func UpdateDocument(doc *entity.Document) error {
	return database.DB.Save(doc).Error
}

// DeleteDocument 删除文档记录
func DeleteDocument(id uint) error {
	return database.DB.Delete(&entity.Document{}, id).Error
}

// ListDocuments 分页查询所有文档列表，支持按标题关键词和状态过滤
func ListDocuments(page, size int, keyword, status string) ([]entity.Document, int64, error) {
	var docs []entity.Document
	var total int64
	query := database.DB.Model(&entity.Document{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&docs).Error
	return docs, total, err
}

// ListDocumentsByUser 分页查询指定用户的文档列表
func ListDocumentsByUser(userID uint, page, size int, keyword, status string) ([]entity.Document, int64, error) {
	var docs []entity.Document
	var total int64
	query := database.DB.Model(&entity.Document{}).Where("user_id = ?", userID)
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&docs).Error
	return docs, total, err
}

// GetAllDocumentsContent 获取所有已解析完成的文档内容（用于本地关键词检索）
func GetAllDocumentsContent() ([]entity.Document, error) {
	var docs []entity.Document
	err := database.DB.Where("content != '' AND status IN ?", []string{"completed", "approved"}).Find(&docs).Error
	return docs, err
}

// CountDocuments 统计文档总数
func CountDocuments() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Document{}).Count(&count).Error
	return count, err
}

// CountDocumentsByStatus 按状态统计文档数量
func CountDocumentsByStatus(status string) (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Document{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// ListDocumentsAdmin 管理员获取文档列表（所有用户的）
func ListDocumentsAdmin(page, size int, keyword string) ([]entity.Document, int64, error) {
	var docs []entity.Document
	var total int64
	query := database.DB.Model(&entity.Document{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&docs).Error
	return docs, total, err
}
