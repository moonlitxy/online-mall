package models

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前的钩子
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate 更新前的钩子
func (base *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
