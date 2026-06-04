package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateAskSession(session *model.AskSession) error {
	return database.DB.Create(session).Error
}

func FindAskSessionByID(id uint) (*model.AskSession, error) {
	var session model.AskSession
	err := database.DB.First(&session, id).Error
	return &session, err
}

func ListAskSessionsByUser(userID uint, page, size int) ([]model.AskSession, int64, error) {
	var sessions []model.AskSession
	var total int64
	query := database.DB.Model(&model.AskSession{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("updated_at DESC").Find(&sessions).Error
	return sessions, total, err
}

func CreateAskMessage(msg *model.AskMessage) error {
	return database.DB.Create(msg).Error
}

func ListAskMessages(sessionID uint, page, size int) ([]model.AskMessage, int64, error) {
	var messages []model.AskMessage
	var total int64
	query := database.DB.Model(&model.AskMessage{}).Where("session_id = ?", sessionID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at ASC").Find(&messages).Error
	return messages, total, err
}

func CountMessagesBySession(sessionID uint) int {
	var count int64
	database.DB.Model(&model.AskMessage{}).Where("session_id = ?", sessionID).Count(&count)
	return int(count)
}

func GetLastMessageBySession(sessionID uint) (*model.AskMessage, error) {
	var msg model.AskMessage
	err := database.DB.Where("session_id = ? AND role = ?", sessionID, "user").Order("created_at DESC").First(&msg).Error
	return &msg, err
}
