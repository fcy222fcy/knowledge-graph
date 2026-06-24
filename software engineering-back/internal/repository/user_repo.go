package repository

import (
	"software_engineering/internal/model/entity"
	"software_engineering/pkg/database"
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

// DeleteUser 删除用户（软删除）
func DeleteUser(id uint) error {
	return database.DB.Delete(&entity.User{}, id).Error
}

// ListUsers 分页获取用户列表
func ListUsers(page, size int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	database.DB.Model(&entity.User{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}

// ListUsersAdmin 管理员获取用户列表（支持搜索）
func ListUsersAdmin(page, size int, keyword string) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := database.DB.Model(&entity.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

// CountUsers 统计用户总数
func CountUsers() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.User{}).Count(&count).Error
	return count, err
}

// CountStudents 统计学生总数（即用户总数）
func CountStudents() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.User{}).Count(&count).Error
	return count, err
}
