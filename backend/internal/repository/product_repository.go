package repository

import (
	"online-mall/internal/models"

	"gorm.io/gorm"
)

// ProductRepository 商品数据访问层
type ProductRepository struct{}

// NewProductRepository 创建商品Repository实例
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// GetByID 根据ID获取商品
func (r *ProductRepository) GetByID(id uint64) (*models.Product, error) {
	var product models.Product
	err := models.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetByIDWithSkus 根据ID获取商品（包含SKU）
func (r *ProductRepository) GetByIDWithSkus(id uint64) (*models.Product, error) {
	var product models.Product
	err := models.DB.Preload("ProductSkus").Preload("Category").Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProducts 分页获取商品列表
func (r *ProductRepository) GetProducts(query *models.ProductQuery) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	db := models.DB.Model(&models.Product{})

	// 分类筛选
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	// 关键字搜索
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}

	// 状态筛选
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	// 热门商品
	if query.IsHot {
		db = db.Where("is_hot = ?", true)
	}

	// 新品
	if query.IsNew {
		db = db.Where("is_new = ?", true)
	}

	// 只查询上架的商品
	if query.Status == 0 {
		db = db.Where("status = ?", 1)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	offset := (query.Page - 1) * query.PageSize

	// 排序
	switch query.Sort {
	case "sales":
		db = db.Order("sales DESC")
	case "price_desc":
		db = db.Order("price DESC")
	case "price_asc":
		db = db.Order("price ASC")
	default:
		db = db.Order("sort DESC, id DESC")
	}

	// 查询列表
	if err := db.Preload("Category").Offset(offset).Limit(query.PageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Create 创建商品
func (r *ProductRepository) Create(product *models.Product) error {
	return models.DB.Create(product).Error
}

// Update 更新商品
func (r *ProductRepository) Update(product *models.Product) error {
	return models.DB.Save(product).Error
}

// Delete 删除商品
func (r *ProductRepository) Delete(id uint64) error {
	return models.DB.Delete(&models.Product{}, id).Error
}

// GetHotProducts 获取热门商品
func (r *ProductRepository) GetHotProducts(limit int) ([]*models.Product, error) {
	var products []*models.Product
	if limit <= 0 {
		limit = 10
	}
	err := models.DB.Where("is_hot = ? AND status = ?", true, 1).
		Order("sales DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// GetNewProducts 获取新品商品
func (r *ProductRepository) GetNewProducts(limit int) ([]*models.Product, error) {
	var products []*models.Product
	if limit <= 0 {
		limit = 10
	}
	err := models.DB.Where("is_new = ? AND status = ?", true, 1).
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// UpdateSales 更新商品销量
func (r *ProductRepository) UpdateSales(productID uint64, quantity int) error {
	return models.DB.Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("sales", gorm.Expr("sales + ?", quantity)).Error
}

// UpdateStock 更新商品库存
func (r *ProductRepository) UpdateStock(productID uint64, quantity int) error {
	return models.DB.Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}
