package repository

import (
	"software_engineering/internal/model/entity"
	"software_engineering/pkg/database"
)

// CreateTeacher 创建教师
func CreateTeacher(teacher *entity.Teacher) error {
	return database.DB.Create(teacher).Error
}

// FindTeacherByUsername 根据用户名查找教师
func FindTeacherByUsername(username string) (*entity.Teacher, error) {
	var teacher entity.Teacher
	err := database.DB.Where("username = ?", username).First(&teacher).Error
	return &teacher, err
}

// FindTeacherByEmail 根据邮箱查找教师
func FindTeacherByEmail(email string) (*entity.Teacher, error) {
	var teacher entity.Teacher
	err := database.DB.Where("email = ?", email).First(&teacher).Error
	return &teacher, err
}

// FindTeacherByID 根据 ID 查找教师
func FindTeacherByID(id uint) (*entity.Teacher, error) {
	var teacher entity.Teacher
	err := database.DB.First(&teacher, id).Error
	return &teacher, err
}

// UpdateTeacher 更新教师信息
func UpdateTeacher(teacher *entity.Teacher) error {
	return database.DB.Save(teacher).Error
}

// DeleteTeacher 删除教师（软删除）
func DeleteTeacher(id uint) error {
	return database.DB.Delete(&entity.Teacher{}, id).Error
}

// ListTeachers 分页获取教师列表
func ListTeachers(page, size int) ([]entity.Teacher, int64, error) {
	var teachers []entity.Teacher
	var total int64
	database.DB.Model(&entity.Teacher{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&teachers).Error
	return teachers, total, err
}

// CountTeachers 统计教师总数
func CountTeachers() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Teacher{}).Count(&count).Error
	return count, err
}
