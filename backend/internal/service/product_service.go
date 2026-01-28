package service

import (
	"errors"
	"online-mall/internal/models"
	"online-mall/internal/repository"
)

// ProductService 商品业务逻辑层
type ProductService struct {
	productRepo *repository.ProductRepository
}

// NewProductService 创建商品Service实例
func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepository(),
	}
}

// GetProduct 获取商品详情
func (s *ProductService) GetProduct(id uint64) (*models.Product, error) {
	return s.productRepo.GetByIDWithSkus(id)
}

// GetProducts 分页获取商品列表
func (s *ProductService) GetProducts(query *models.ProductQuery) ([]*models.Product, int64, error) {
	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	return s.productRepo.GetProducts(query)
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(product *models.Product) error {
	// 验证分类是否存在
	var category models.Category
	if err := models.DB.First(&category, product.CategoryID).Error; err != nil {
		return errors.New("分类不存在")
	}

	// 设置默认值
	if product.Status == 0 {
		product.Status = 1
	}

	return s.productRepo.Create(product)
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(product *models.Product) error {
	// 检查商品是否存在
	_, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return errors.New("商品不存在")
	}

	// 验证分类是否存在（如果修改了分类）
	if product.CategoryID > 0 {
		var category models.Category
		if err := models.DB.First(&category, product.CategoryID).Error; err != nil {
			return errors.New("分类不存在")
		}
	}

	return s.productRepo.Update(product)
}

// DeleteProduct 删除商品
func (s *ProductService) DeleteProduct(id uint64) error {
	// 检查商品是否存在
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return errors.New("商品不存在")
	}

	// 检查是否可以删除（可以根据业务需求添加更多检查）
	// 例如：检查是否有未完成的订单等

	// 软删除
	return s.productRepo.Delete(id)
}

// GetHotProducts 获取热门商品
func (s *ProductService) GetHotProducts(limit int) ([]*models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.productRepo.GetHotProducts(limit)
}

// GetNewProducts 获取新品商品
func (s *ProductService) GetNewProducts(limit int) ([]*models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.productRepo.GetNewProducts(limit)
}

// UpdateSales 更新商品销量
func (s *ProductService) UpdateSales(productID uint64, quantity int) error {
	return s.productRepo.UpdateSales(productID, quantity)
}

// UpdateStock 更新商品库存
func (s *ProductService) UpdateStock(productID uint64, quantity int) error {
	// 检查商品是否存在
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return errors.New("商品不存在")
	}

	// 检查库存是否充足
	if quantity < 0 && product.Stock < -quantity {
		return errors.New("库存不足")
	}

	return s.productRepo.UpdateStock(productID, quantity)
}

// GetCategoryProducts 获取分类下的商品
func (s *ProductService) GetCategoryProducts(categoryID uint64, page, pageSize int) ([]*models.Product, int64, error) {
	query := &models.ProductQuery{
		CategoryID: categoryID,
		Page:       page,
		PageSize:   pageSize,
		Status:     1, // 只查询上架商品
	}
	return s.productRepo.GetProducts(query)
}

// SearchProducts 搜索商品
func (s *ProductService) SearchProducts(keyword string, page, pageSize int) ([]*models.Product, int64, error) {
	query := &models.ProductQuery{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Status:   1, // 只查询上架商品
	}
	return s.productRepo.GetProducts(query)
}
