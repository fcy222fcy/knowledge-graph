package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

// CreateUser 创建用户
func CreateUser(user *entity.User) error {
	return database.DB.Create(user).Error
}

// FindUserByUsername 根据用户名查找用户
func FindUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindUserByEmail 根据邮箱查找用户
func FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindUserByID 根据 ID 查找用户
func FindUserByID(id uint) (*entity.User, error) {
	var user entity.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

// UpdateUser 更新用户信息
func UpdateUser(user *entity.User) error {
	return database.DB.Save(user).Error
}

// ListUsers 分页获取用户列表
func ListUsers(page, size int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	database.DB.Model(&entity.User{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
