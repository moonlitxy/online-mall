package repository

import (
	"online-mall/internal/models"

	"gorm.io/gorm"
)

// CategoryRepository 分类数据访问层
type CategoryRepository struct{}

// NewCategoryRepository 创建分类Repository实例
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(id uint64) (*models.Category, error) {
	var category models.Category
	err := models.DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有分类
func (r *CategoryRepository) GetAll() ([]*models.Category, error) {
	var categories []*models.Category
	err := models.DB.Order("sort ASC, id ASC").Find(&categories).Error
	return categories, err
}

// GetByParentID 根据父分类ID获取子分类
func (r *CategoryRepository) GetByParentID(parentID uint64) ([]*models.Category, error) {
	var categories []*models.Category
	err := models.DB.Where("parent_id = ?", parentID).
		Order("sort ASC, id ASC").
		Find(&categories).Error
	return categories, err
}

// GetTree 获取分类树
func (r *CategoryRepository) GetTree() ([]*models.Category, error) {
	// 获取所有分类
	var categories []*models.Category
	if err := models.DB.Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 构建树形结构
	return r.buildTree(categories, 0), nil
}

// buildTree 构建树形结构
func (r *CategoryRepository) buildTree(categories []*models.Category, parentID uint64) []*models.Category {
	var tree []*models.Category
	for _, category := range categories {
		if category.ParentID == parentID {
			children := r.buildTree(categories, category.ID)
			if len(children) > 0 {
				// 这里可以将children添加到category中，但由于Category模型没有Children字段，
				// 我们可以在Service层处理或者修改模型
			}
			tree = append(tree, category)
		}
	}
	return tree
}

// Create 创建分类
func (r *CategoryRepository) Create(category *models.Category) error {
	return models.DB.Create(category).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(category *models.Category) error {
	return models.DB.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint64) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否有子分类
		var count int64
		if err := tx.Model(&models.Category{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return gorm.ErrDuplicatedKey // 使用这个错误表示有子分类
		}

		// 检查是否有商品使用该分类
		var productCount int64
		if err := tx.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
			return err
		}
		if productCount > 0 {
			return gorm.ErrDuplicatedKey // 使用这个错误表示有商品
		}

		// 删除分类
		return tx.Delete(&models.Category{}, id).Error
	})
}

// UpdateStatus 更新分类状态
func (r *CategoryRepository) UpdateStatus(id uint64, status int) error {
	return models.DB.Model(&models.Category{}).
		Where("id = ?", id).
		Update("status", status).Error
}
