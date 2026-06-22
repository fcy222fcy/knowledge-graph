package repository

import (
	"software_engineering/internal/model/entity"

	"gorm.io/gorm"
)

// UserRepository 定义用户仓库接口
type UserRepository interface {
	Create(user *entity.User) error                        // 创建用户
	FindByUsername(username string) (*entity.User, error)   // 根据用户名查找
	FindByEmail(email string) (*entity.User, error)        // 根据邮箱查找
	FindByID(id uint) (*entity.User, error)                // 根据ID查找
	Update(user *entity.User) error                        // 更新用户
	List(page, size int) ([]entity.User, int64, error)     // 分页列表
}

// GormUserRepository GORM 用户仓库实现
type GormUserRepository struct {
	db *gorm.DB // 数据库连接
}

// NewGormUserRepository 创建 GORM 用户仓库实例
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Create 创建用户
func (r *GormUserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindByUsername 根据用户名查找用户
func (r *GormUserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindByEmail 根据邮箱查找用户
func (r *GormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindByID 根据 ID 查找用户
func (r *GormUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// Update 更新用户信息
func (r *GormUserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// List 分页获取用户列表
func (r *GormUserRepository) List(page, size int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	r.db.Model(&entity.User{}).Count(&total)
	err := r.db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
