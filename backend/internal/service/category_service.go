package service

import (
	"errors"
	"online-mall/internal/models"
	"online-mall/internal/repository"
)

// CategoryService 分类业务逻辑层
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryService 创建分类Service实例
func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewCategoryRepository(),
	}
}

// GetCategory 获取分类详情
func (s *CategoryService) GetCategory(id uint64) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

// GetCategories 获取所有分类
func (s *CategoryService) GetCategories() ([]*models.Category, error) {
	return s.categoryRepo.GetAll()
}

// GetCategoryTree 获取分类树
func (s *CategoryService) GetCategoryTree() ([]*models.Category, error) {
	return s.categoryRepo.GetTree()
}

// GetSubCategories 获取子分类
func (s *CategoryService) GetSubCategories(parentID uint64) ([]*models.Category, error) {
	return s.categoryRepo.GetByParentID(parentID)
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(category *models.Category) error {
	// 验证父分类是否存在（如果不是顶级分类）
	if category.ParentID > 0 {
		parent, err := s.categoryRepo.GetByID(category.ParentID)
		if err != nil {
			return errors.New("父分类不存在")
		}
		// 设置层级
		category.Level = parent.Level + 1
	} else {
		category.Level = 1
	}

	// 设置默认值
	if category.Status == 0 {
		category.Status = 1
	}

	return s.categoryRepo.Create(category)
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(category *models.Category) error {
	// 检查分类是否存在
	_, err := s.categoryRepo.GetByID(category.ID)
	if err != nil {
		return errors.New("分类不存在")
	}

	// 如果修改了父分类，需要更新层级
	if category.ParentID > 0 {
		parent, err := s.categoryRepo.GetByID(category.ParentID)
		if err != nil {
			return errors.New("父分类不存在")
		}
		// 不能将分类设置为自己的子分类
		if category.ID == parent.ID {
			return errors.New("不能将分类设置为自己的子分类")
		}
		category.Level = parent.Level + 1
	} else {
		category.Level = 1
	}

	return s.categoryRepo.Update(category)
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint64) error {
	// 检查分类是否存在
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	// 删除分类（Repository中会检查是否有子分类和商品）
	err = s.categoryRepo.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New("分类不存在")
		}
		return errors.New("删除失败，该分类下有子分类或商品")
	}

	return nil
}

// UpdateCategoryStatus 更新分类状态
func (s *CategoryService) UpdateCategoryStatus(id uint64, status int) error {
	// 检查分类是否存在
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	return s.categoryRepo.UpdateStatus(id, status)
}

// GetProductsCount 获取分类下的商品数量
func (s *CategoryService) GetProductsCount(categoryID uint64) (int64, error) {
	var count int64
	err := models.DB.Model(&models.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}
