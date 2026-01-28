package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	BaseModel
	Username    string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=6"`
	Phone       string     `gorm:"type:varchar(20);uniqueIndex" json:"phone" validate:"e164"`
	Email       string     `gorm:"type:varchar(100);uniqueIndex" json:"email" validate:"email"`
	Nickname    string     `gorm:"type:varchar(50)" json:"nickname"`
	Avatar      string     `gorm:"type:varchar(255)" json:"avatar"`
	Status      int        `gorm:"type:tinyint;default:1" json:"status"` // 1-正常，0-禁用
	LastLoginAt *time.Time `json:"last_login_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 加密密码
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 如果更新了密码，需要加密
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserQuery 用户查询结构体
type UserQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Username string `form:"username" json:"username"`
	Phone    string `form:"phone" json:"phone"`
	Email    string `form:"email" json:"email"`
	Status   int    `form:"status" json:"status"`
}
