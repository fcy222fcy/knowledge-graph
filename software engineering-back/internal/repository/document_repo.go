package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateDocument(doc *model.Document) error {
	return database.DB.Create(doc).Error
}

func FindDocumentByID(id uint) (*model.Document, error) {
	var doc model.Document
	err := database.DB.First(&doc, id).Error
	return &doc, err
}

func UpdateDocument(doc *model.Document) error {
	return database.DB.Save(doc).Error
}

func DeleteDocument(id uint) error {
	return database.DB.Delete(&model.Document{}, id).Error
}

func ListDocuments(page, size int, keyword, status string) ([]model.Document, int64, error) {
	var docs []model.Document
	var total int64
	query := database.DB.Model(&model.Document{})
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
