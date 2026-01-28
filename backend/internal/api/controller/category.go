package controller

import (
	"fmt"
	"online-mall/internal/models"
	"online-mall/internal/service"
	"online-mall/internal/utils"

	"github.com/gin-gonic/gin"
)

// CategoryService 分类服务实例
var categoryService = service.NewCategoryService()

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint64 `json:"parent_id"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" binding:"oneof=0 1"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=1"`
	ParentID *uint64 `json:"parent_id" binding:"omitempty,min=0"`
	Sort     *int    `json:"sort"`
	Status   *int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// GetCategories 获取所有分类
func GetCategories(c *gin.Context) {
	categories, err := categoryService.GetCategories()
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, categories)
}

// GetCategoryTree 获取分类树
func GetCategoryTree(c *gin.Context) {
	categories, err := categoryService.GetCategoryTree()
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, categories)
}

// GetCategory 获取分类详情
func GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "分类ID不能为空")
		return
	}

	var categoryID uint64
	if _, err := fmt.Sscanf(id, "%d", &categoryID); err != nil {
		utils.ParamError(c, "分类ID格式错误")
		return
	}

	category, err := categoryService.GetCategory(categoryID)
	if err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	utils.Success(c, category)
}

// GetSubCategories 获取子分类
func GetSubCategories(c *gin.Context) {
	parentID := c.Param("id")
	if parentID == "" {
		utils.ParamError(c, "父分类ID不能为空")
		return
	}

	var pid uint64
	if _, err := fmt.Sscanf(parentID, "%d", &pid); err != nil {
		utils.ParamError(c, "分类ID格式错误")
		return
	}

	categories, err := categoryService.GetSubCategories(pid)
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, categories)
}

// CreateCategory 创建分类（管理员）
func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	category := &models.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   req.Status,
	}

	if err := categoryService.CreateCategory(category); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, category)
}

// UpdateCategory 更新分类（管理员）
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "分类ID不能为空")
		return
	}

	var categoryID uint64
	if _, err := fmt.Sscanf(id, "%d", &categoryID); err != nil {
		utils.ParamError(c, "分类ID格式错误")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	category, err := categoryService.GetCategory(categoryID)
	if err != nil {
		utils.NotFound(c, "分类不存在")
		return
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.ParentID != nil {
		category.ParentID = *req.ParentID
	}
	if req.Sort != nil {
		category.Sort = *req.Sort
	}
	if req.Status != nil {
		category.Status = *req.Status
	}

	if err := categoryService.UpdateCategory(category); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Updated(c, category)
}

// DeleteCategory 删除分类（管理员）
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "分类ID不能为空")
		return
	}

	var categoryID uint64
	if _, err := fmt.Sscanf(id, "%d", &categoryID); err != nil {
		utils.ParamError(c, "分类ID格式错误")
		return
	}

	if err := categoryService.DeleteCategory(categoryID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, map[string]string{
		"message": "删除成功",
	})
}

// UpdateCategoryStatus 更新分类状态（管理员）
func UpdateCategoryStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "分类ID不能为空")
		return
	}

	var categoryID uint64
	if _, err := fmt.Sscanf(id, "%d", &categoryID); err != nil {
		utils.ParamError(c, "分类ID格式错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	if err := categoryService.UpdateCategoryStatus(categoryID, req.Status); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, map[string]string{
		"message": "状态更新成功",
	})
}
