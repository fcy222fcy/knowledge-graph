package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindUserByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func ListUsers(page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	database.DB.Model(&model.User{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
