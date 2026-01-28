package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Order 订单模型
type Order struct {
	BaseModel
	OrderNo        string      `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"`
	UserID         uint64      `gorm:"not null;index" json:"user_id"`
	AddressID      uint64      `gorm:"not null" json:"address_id"`
	TotalAmount    float64     `gorm:"type:decimal(10,2);not null" json:"total_amount" validate:"required,gte=0"`
	Freight        float64     `gorm:"type:decimal(10,2);default:0.00" json:"freight" validate:"gte=0"`
	DiscountAmount float64     `gorm:"type:decimal(10,2);default:0.00" json:"discount_amount" validate:"gte=0"`
	PayAmount      float64     `gorm:"type:decimal(10,2);not null" json:"pay_amount" validate:"required,gte=0"`
	PayStatus      int         `gorm:"type:tinyint;default:0" json:"pay_status"` // 0-未支付，1-已支付
	PayTime        *time.Time  `json:"pay_time"`
	PaymentMethod  string      `gorm:"type:varchar(20)" json:"payment_method"`
	OrderStatus    int         `gorm:"type:tinyint;default:0" json:"order_status"` // 0-待付款，1-待发货，2-待收货，3-已完成，4-已取消
	CancelReason   string      `gorm:"type:varchar(255)" json:"cancel_reason"`
	CancelTime     *time.Time  `json:"cancel_time"`
	Remark         string      `gorm:"type:varchar(255)" json:"remark"`
	User           User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Address        Address     `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	OrderItems     []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

// TableName 表名
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate 创建前钩子
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	// 生成订单号
	if o.OrderNo == "" {
		o.OrderNo = o.generateOrderNo()
	}
	return nil
}

// generateOrderNo 生成订单号
func (o *Order) generateOrderNo() string {
	timestamp := time.Now().Format("20060102150405")
	random := rand.Intn(10000)
	return fmt.Sprintf("%s%04d", timestamp, random)
}

// OrderItem 订单商品模型
type OrderItem struct {
	BaseModel
	OrderID        uint64  `gorm:"not null;index" json:"order_id"`
	ProductID      uint64  `gorm:"not null;index" json:"product_id"`
	SKUID          uint64  `gorm:"not null;index" json:"sku_id"`
	ProductName    string  `gorm:"type:varchar(255);not null" json:"product_name"`
	ProductImage   string  `gorm:"type:varchar(255)" json:"product_image"`
	Specifications string  `gorm:"type:text" json:"specifications"` // JSON格式存储规格信息
	Price          float64 `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	Quantity       int     `gorm:"not null" json:"quantity" validate:"required,gte=1"`
	TotalAmount    float64 `gorm:"type:decimal(10,2);not null" json:"total_amount" validate:"required,gte=0"`
	Order          Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product        Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 表名
func (OrderItem) TableName() string {
	return "order_items"
}

// GetSpecifications 获取规格信息
func (oi *OrderItem) GetSpecifications() map[string]string {
	var specs map[string]string
	if oi.Specifications != "" {
		_ = json.Unmarshal([]byte(oi.Specifications), &specs)
	}
	return specs
}

// SetSpecifications 设置规格信息
func (oi *OrderItem) SetSpecifications(specs map[string]string) {
	data, _ := json.Marshal(specs)
	oi.Specifications = string(data)
}

// CartItem 购物车模型
type CartItem struct {
	BaseModel
	UserID     uint64     `gorm:"not null;index" json:"user_id"`
	ProductID  uint64     `gorm:"not null;index" json:"product_id"`
	SKUID      uint64     `gorm:"not null;index" json:"sku_id"`
	Quantity   int        `gorm:"not null;default:1" json:"quantity" validate:"required,gte=1"`
	Selected   bool       `gorm:"type:boolean;default:true" json:"selected"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product    Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductSKU ProductSKU `gorm:"foreignKey:SKUID" json:"product_sku,omitempty"`
}

// TableName 表名
func (CartItem) TableName() string {
	return "cart_items"
}

// OrderQuery 订单查询结构体
type OrderQuery struct {
	Status    int    `form:"status" json:"status"`
	Page      int    `form:"page" json:"page"`
	PageSize  int    `form:"page_size" json:"page_size"`
	OrderNo   string `form:"order_no" json:"order_no"`
	UserID    uint64 `form:"user_id" json:"user_id"`
	StartTime string `form:"start_time" json:"start_time"`
	EndTime   string `form:"end_time" json:"end_time"`
}

// CartQuery 购物车查询结构体
type CartQuery struct {
	UserID   uint64 `form:"user_id" json:"user_id"`
	Selected *bool  `form:"selected" json:"selected"` // 可选，true-选中，false-未选中
}
