package repository

import (
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
)

func CreateUser(user *entity.User) error {
	return database.DB.Create(user).Error
}

func FindUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func FindUserByID(id uint) (*entity.User, error) {
	var user entity.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *entity.User) error {
	return database.DB.Save(user).Error
}

func ListUsers(page, size int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	database.DB.Model(&entity.User{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
