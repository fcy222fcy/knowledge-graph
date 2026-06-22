package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含 ID、创建时间、更新时间和软删除字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey;comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间" json:"-"` // 软删除，查询时自动过滤
}
