package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Coupon 优惠券模型
type Coupon struct {
	BaseModel
	Name      string  `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Type      int     `gorm:"type:tinyint;not null" json:"type"` // 1-满减券，2-折扣券
	Value     float64 `gorm:"type:decimal(10,2);not null" json:"value" validate:"required,gte=0"`
	MinAmount float64 `gorm:"type:decimal(10,2);default:0.00" json:"min_amount" validate:"gte=0"`
	StartTime string  `gorm:"type:datetime;not null" json:"start_time" validate:"required"`
	EndTime   string  `gorm:"type:datetime;not null" json:"end_time" validate:"required"`
	Stock     int     `gorm:"default:0" json:"stock" validate:"gte=0"`
	UsedCount int     `gorm:"default:0" json:"used_count"`          // 已使用数量
	Status    int     `gorm:"type:tinyint;default:1" json:"status"` // 1-可用，0-停用
}

// TableName 表名
func (Coupon) TableName() string {
	return "coupons"
}

// UserCoupon 用户优惠券模型
type UserCoupon struct {
	BaseModel
	UserID   uint64     `gorm:"not null;index" json:"user_id"`
	CouponID uint64     `gorm:"not null;index" json:"coupon_id"`
	OrderID  *uint64    `gorm:"index" json:"order_id"`                // 使用的订单ID
	Status   int        `gorm:"type:tinyint;default:0" json:"status"` // 0-未使用，1-已使用，2-已过期
	UsedTime *time.Time `json:"used_time"`
	User     User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coupon   Coupon     `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

// TableName 表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}

// CouponQuery 优惠券查询结构体
type CouponQuery struct {
	Page       int    `form:"page" json:"page"`
	PageSize   int    `form:"page_size" json:"page_size"`
	Status     int    `form:"status" json:"status"` // 1-可用，0-停用
	Type       int    `form:"type" json:"type"`     // 1-满减券，2-折扣券
	UserID     uint64 `form:"user_id" json:"user_id"`
	StatusType int    `form:"status_type" json:"status_type"` // 0-未使用，1-已使用，2-已过期
}

// GetUserCouponKey 获取用户优惠券缓存key
func (uc *UserCoupon) GetUserCouponKey() string {
	return fmt.Sprintf("user:coupon:%d:%d", uc.UserID, uc.CouponID)
}

// GetCouponKey 获取优惠券缓存key
func (c *Coupon) GetCouponKey() string {
	return fmt.Sprintf("coupon:%d", c.ID)
}

// CheckIsValid 检查优惠券是否有效
func (c *Coupon) CheckIsValid() bool {
	// 检查状态
	if c.Status != 1 {
		return false
	}

	// 检查库存
	if c.Stock > 0 && c.UsedCount >= c.Stock {
		return false
	}

	// 检查时间
	now := time.Now()
	startTime, _ := time.Parse("2006-01-02 15:04:05", c.StartTime)
	endTime, _ := time.Parse("2006-01-02 15:04:05", c.EndTime)

	if now.Before(startTime) || now.After(endTime) {
		return false
	}

	return true
}

// GetDiscountAmount 计算优惠金额
func (c *Coupon) GetDiscountAmount(totalAmount float64) float64 {
	if !c.CheckIsValid() {
		return 0
	}

	if totalAmount < c.MinAmount {
		return 0
	}

	switch c.Type {
	case 1: // 满减券
		return c.Value
	case 2: // 折扣券
		if c.Value >= 1 || c.Value <= 0 {
			return 0
		}
		discount := totalAmount * (1 - c.Value)
		return discount
	default:
		return 0
	}
}

// GetDiscountDisplay 获取优惠显示文本
func (c *Coupon) GetDiscountDisplay() string {
	switch c.Type {
	case 1: // 满减券
		return fmt.Sprintf("满%.2f减%.2f", c.MinAmount, c.Value)
	case 2: // 折扣券
		return fmt.Sprintf("%.0f折", c.Value*10)
	default:
		return ""
	}
}

// UpdateStatus 更新优惠券状态
func (c *Coupon) UpdateStatus(db *gorm.DB, newStatus int) error {
	c.Status = newStatus
	return db.Save(c).Error
}

// UpdateUsedCount 更新已使用数量
func (c *Coupon) UpdateUsedCount(db *gorm.DB, increment int) error {
	return db.Model(c).UpdateColumn("used_count", gorm.Expr("used_count + ?", increment)).Error
}
