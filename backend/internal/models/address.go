package models

import (
	"gorm.io/gorm"
)

// Address 收货地址模型
type Address struct {
	BaseModel
	UserID    uint64 `gorm:"not null;index" json:"user_id"`
	Name      string `gorm:"type:varchar(50);not null" json:"name" validate:"required"`
	Phone     string `gorm:"type:varchar(20);not null" json:"phone" validate:"required,e164"`
	Province  string `gorm:"type:varchar(50);not null" json:"province" validate:"required"`
	City      string `gorm:"type:varchar(50);not null" json:"city" validate:"required"`
	District  string `gorm:"type:varchar(50);not null" json:"district" validate:"required"`
	Detail    string `gorm:"type:varchar(255);not null" json:"detail" validate:"required"`
	Postcode  string `gorm:"type:varchar(10)" json:"postcode"`
	Tag       string `gorm:"type:varchar(20)" json:"tag"` // 家、公司、学校等
	IsDefault bool   `gorm:"type:boolean;default:false" json:"is_default"`
}

// TableName 表名
func (Address) TableName() string {
	return "addresses"
}

// BeforeCreate 创建前钩子
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	// 如果设置为默认地址，需要取消其他默认地址
	if a.IsDefault {
		tx.Model(&Address{}).Where("user_id = ?", a.UserID).Update("is_default", false)
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (a *Address) BeforeUpdate(tx *gorm.DB) error {
	// 如果更新了默认地址状态
	if tx.Statement.Changed("IsDefault") && a.IsDefault {
		tx.Model(&Address{}).Where("user_id = ?", a.UserID).Update("is_default", false)
	}
	return nil
}

// AddressQuery 地址查询结构体
type AddressQuery struct {
	UserID uint64 `form:"user_id" json:"user_id"` // 用户ID
}
