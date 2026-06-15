package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

func CreateAskSession(session *entity.AskSession) error {
	return database.DB.Create(session).Error
}

func FindAskSessionByID(id uint) (*entity.AskSession, error) {
	var session entity.AskSession
	err := database.DB.First(&session, id).Error
	return &session, err
}

func ListAskSessionsByUser(userID uint, page, size int) ([]entity.AskSession, int64, error) {
	var sessions []entity.AskSession
	var total int64
	query := database.DB.Model(&entity.AskSession{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("updated_at DESC").Find(&sessions).Error
	return sessions, total, err
}

func CreateAskMessage(msg *entity.AskMessage) error {
	return database.DB.Create(msg).Error
}

func ListAskMessages(sessionID uint, page, size int) ([]entity.AskMessage, int64, error) {
	var messages []entity.AskMessage
	var total int64
	query := database.DB.Model(&entity.AskMessage{}).Where("session_id = ?", sessionID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at ASC").Find(&messages).Error
	return messages, total, err
}

func CountMessagesBySession(sessionID uint) int {
	var count int64
	database.DB.Model(&entity.AskMessage{}).Where("session_id = ?", sessionID).Count(&count)
	return int(count)
}

func GetLastMessageBySession(sessionID uint) (*entity.AskMessage, error) {
	var msg entity.AskMessage
	err := database.DB.Where("session_id = ? AND role = ?", sessionID, "user").Order("created_at DESC").First(&msg).Error
	return &msg, err
}

func ListRecentMessages(sessionID uint, limit int) ([]entity.AskMessage, error) {
	var messages []entity.AskMessage
	err := database.DB.Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	// 反转为时间正序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, err
}
