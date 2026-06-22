package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CreateAskSession 创建问答会话
func CreateAskSession(session *entity.AskSession) error {
	return database.DB.Create(session).Error
}

// FindAskSessionByID 根据 ID 查找问答会话
func FindAskSessionByID(id uint) (*entity.AskSession, error) {
	var session entity.AskSession
	err := database.DB.First(&session, id).Error
	return &session, err
}

// ListAskSessionsByUser 分页获取用户的问答会话列表
func ListAskSessionsByUser(userID uint, page, size int) ([]entity.AskSession, int64, error) {
	var sessions []entity.AskSession
	var total int64
	query := database.DB.Model(&entity.AskSession{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("updated_at DESC").Find(&sessions).Error
	return sessions, total, err
}

// CreateAskMessage 创建问答消息
func CreateAskMessage(msg *entity.AskMessage) error {
	return database.DB.Create(msg).Error
}

// ListAskMessages 分页获取指定会话的消息记录
func ListAskMessages(sessionID uint, page, size int) ([]entity.AskMessage, int64, error) {
	var messages []entity.AskMessage
	var total int64
	query := database.DB.Model(&entity.AskMessage{}).Where("session_id = ?", sessionID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at ASC").Find(&messages).Error
	return messages, total, err
}

// CountMessagesBySession 统计指定会话的消息数量
func CountMessagesBySession(sessionID uint) int {
	var count int64
	database.DB.Model(&entity.AskMessage{}).Where("session_id = ?", sessionID).Count(&count)
	return int(count)
}

// GetLastMessageBySession 获取指定会话的最后一条用户消息
func GetLastMessageBySession(sessionID uint) (*entity.AskMessage, error) {
	var msg entity.AskMessage
	err := database.DB.Where("session_id = ? AND role = ?", sessionID, "user").Order("created_at DESC").First(&msg).Error
	return &msg, err
}

// ListRecentMessages 获取指定会话的最近 N 条消息，按时间正序返回
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
