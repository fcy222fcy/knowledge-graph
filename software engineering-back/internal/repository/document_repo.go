package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

func CreateDocument(doc *entity.Document) error {
	return database.DB.Create(doc).Error
}

func FindDocumentByID(id uint) (*entity.Document, error) {
	var doc entity.Document
	err := database.DB.First(&doc, id).Error
	return &doc, err
}

func UpdateDocument(doc *entity.Document) error {
	return database.DB.Save(doc).Error
}

func DeleteDocument(id uint) error {
	return database.DB.Delete(&entity.Document{}, id).Error
}

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

func GetAllDocumentsContent() ([]entity.Document, error) {
	var docs []entity.Document
	err := database.DB.Where("content != '' AND status = ?", "completed").Find(&docs).Error
	return docs, err
}
