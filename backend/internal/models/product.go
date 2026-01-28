package models

import (
	"encoding/json"
)

// Category 商品分类模型
type Category struct {
	BaseModel
	Name     string `gorm:"type:varchar(50);not null" json:"name" validate:"required"`
	ParentID uint64 `gorm:"default:0;index" json:"parent_id"`     // 0表示顶级分类
	Level    int    `gorm:"default:1" json:"level"`               // 层级
	Sort     int    `gorm:"default:0" json:"sort"`                // 排序
	Status   int    `gorm:"type:tinyint;default:1" json:"status"` // 1-显示，0-隐藏
}

// TableName 表名
func (Category) TableName() string {
	return "categories"
}

// Product 商品模型
type Product struct {
	BaseModel
	Name          string       `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	CategoryID    uint64       `gorm:"not null;index" json:"category_id" validate:"required"`
	Description   string       `gorm:"type:text" json:"description"`
	Price         float64      `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	OriginalPrice *float64     `gorm:"type:decimal(10,2)" json:"original_price"`
	Stock         int          `gorm:"default:0" json:"stock" validate:"gte=0"`
	Sales         int          `gorm:"default:0" json:"sales"`  // 销量
	Images        string       `gorm:"type:text" json:"images"` // JSON格式存储图片数组
	VideoURL      string       `gorm:"type:varchar(255)" json:"video_url"`
	Status        int          `gorm:"type:tinyint;default:1" json:"status"` // 1-上架，0-下架
	IsHot         bool         `gorm:"type:boolean;default:false" json:"is_hot"`
	IsNew         bool         `gorm:"type:boolean;default:false" json:"is_new"`
	Sort          int          `gorm:"default:0" json:"sort"`
	Category      Category     `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ProductSkus   []ProductSKU `gorm:"foreignKey:ProductID" json:"product_skus,omitempty"`
}

// TableName 表名
func (Product) TableName() string {
	return "products"
}

// GetImages 获取商品图片数组
func (p *Product) GetImages() []string {
	var images []string
	if p.Images != "" {
		_ = json.Unmarshal([]byte(p.Images), &images)
	}
	return images
}

// SetImages 设置商品图片数组
func (p *Product) SetImages(images []string) {
	data, _ := json.Marshal(images)
	p.Images = string(data)
}

// ProductSKU 商品SKU模型
type ProductSKU struct {
	BaseModel
	ProductID      uint64  `gorm:"not null;index" json:"product_id"`
	Name           string  `gorm:"type:varchar(255);not null" json:"name"`
	Specifications string  `gorm:"type:text;not null" json:"specifications"` // JSON格式存储规格信息
	Price          float64 `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,gte=0"`
	Stock          int     `gorm:"default:0" json:"stock" validate:"gte=0"`
	Sales          int     `gorm:"default:0" json:"sales"` // SKU销量
	Image          string  `gorm:"type:varchar(255)" json:"image"`
	Product        Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 表名
func (ProductSKU) TableName() string {
	return "product_skus"
}

// GetSpecifications 获取规格信息
func (s *ProductSKU) GetSpecifications() map[string]string {
	var specs map[string]string
	if s.Specifications != "" {
		_ = json.Unmarshal([]byte(s.Specifications), &specs)
	}
	return specs
}

// SetSpecifications 设置规格信息
func (s *ProductSKU) SetSpecifications(specs map[string]string) {
	data, _ := json.Marshal(specs)
	s.Specifications = string(data)
}

// ProductQuery 商品查询结构体
type ProductQuery struct {
	Page       int    `form:"page" json:"page"`
	PageSize   int    `form:"page_size" json:"page_size"`
	CategoryID uint64 `form:"category_id" json:"category_id"`
	Keyword    string `form:"keyword" json:"keyword"`
	Sort       string `form:"sort" json:"sort"` // sales, price_desc, price_asc
	Status     int    `form:"status" json:"status"`
	IsHot      bool   `form:"is_hot" json:"is_hot"`
	IsNew      bool   `form:"is_new" json:"is_new"`
}
