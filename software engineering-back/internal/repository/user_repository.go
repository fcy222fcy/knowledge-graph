package repository

import (
	"software_engineering/internal/model/entity"

	"gorm.io/gorm"
)

// UserRepository 定义用户仓库接口
type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	Update(user *entity.User) error
	List(page, size int) ([]entity.User, int64, error)
}

// GormUserRepository GORM 用户仓库实现
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建 GORM 用户仓库实例
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *GormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *GormUserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *GormUserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) List(page, size int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	r.db.Model(&entity.User{}).Count(&total)
	err := r.db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
