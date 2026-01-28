package controller

import (
	"fmt"
	"online-mall/internal/models"
	"online-mall/internal/service"
	"online-mall/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductService 商品服务实例
var productService = service.NewProductService()

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name          string   `json:"name" binding:"required"`
	CategoryID    uint64   `json:"category_id" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required,gte=0"`
	OriginalPrice *float64 `json:"original_price" binding:"omitempty,gte=0"`
	Stock         int      `json:"stock" binding:"gte=0"`
	Images        []string `json:"images"`
	VideoURL      string   `json:"video_url"`
	Status        int      `json:"status" binding:"oneof=0 1"`
	IsHot         bool     `json:"is_hot"`
	IsNew         bool     `json:"is_new"`
	Sort          int      `json:"sort"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	Name          *string  `json:"name" binding:"omitempty,min=1"`
	CategoryID    *uint64  `json:"category_id" binding:"omitempty,min=1"`
	Description   *string  `json:"description"`
	Price         *float64 `json:"price" binding:"omitempty,gte=0"`
	OriginalPrice *float64 `json:"original_price" binding:"omitempty,gte=0"`
	Stock         *int     `json:"stock" binding:"omitempty,gte=0"`
	Images        []string `json:"images"`
	VideoURL      *string  `json:"video_url"`
	Status        *int     `json:"status" binding:"omitempty,oneof=0 1"`
	IsHot         *bool    `json:"is_hot"`
	IsNew         *bool    `json:"is_new"`
	Sort          *int     `json:"sort"`
}

// ProductQuery 商品查询请求
type ProductQuery struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	CategoryID uint64 `form:"category_id" binding:"omitempty,min=1"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort" binding:"omitempty,oneof=sales price_desc price_asc"`
	IsHot      bool   `form:"is_hot"`
	IsNew      bool   `form:"is_new"`
}

// GetProducts 获取商品列表
func GetProducts(c *gin.Context) {
	var query ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 转换为模型查询
	productQuery := &models.ProductQuery{
		Page:       query.Page,
		PageSize:   query.PageSize,
		CategoryID: query.CategoryID,
		Keyword:    query.Keyword,
		Sort:       query.Sort,
		IsHot:      query.IsHot,
		IsNew:      query.IsNew,
	}

	// 调用服务层
	products, total, err := productService.GetProducts(productQuery)
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, map[string]interface{}{
		"list":      products,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

// GetProduct 获取商品详情
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "商品ID不能为空")
		return
	}

	var productID uint64
	if _, err := fmt.Sscanf(id, "%d", &productID); err != nil {
		utils.ParamError(c, "商品ID格式错误")
		return
	}

	product, err := productService.GetProduct(productID)
	if err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	utils.Success(c, product)
}

// GetProductSkus 获取商品SKU列表
func GetProductSkus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "商品ID不能为空")
		return
	}

	var productID uint64
	if _, err := fmt.Sscanf(id, "%d", &productID); err != nil {
		utils.ParamError(c, "商品ID格式错误")
		return
	}

	product, err := productService.GetProduct(productID)
	if err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	utils.Success(c, product.ProductSkus)
}

// GetHotProducts 获取热门商品
func GetHotProducts(c *gin.Context) {
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	products, err := productService.GetHotProducts(limit)
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, products)
}

// GetNewProducts 获取新品商品
func GetNewProducts(c *gin.Context) {
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	products, err := productService.GetNewProducts(limit)
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, products)
}

// CreateProduct 创建商品（管理员）
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 转换为模型
	product := &models.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Status:        req.Status,
		IsHot:         req.IsHot,
		IsNew:         req.IsNew,
		Sort:          req.Sort,
		VideoURL:      req.VideoURL,
	}

	// 设置图片
	if len(req.Images) > 0 {
		product.SetImages(req.Images)
	}

	// 调用服务层
	if err := productService.CreateProduct(product); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, product)
}

// UpdateProduct 更新商品（管理员）
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "商品ID不能为空")
		return
	}

	var productID uint64
	if _, err := fmt.Sscanf(id, "%d", &productID); err != nil {
		utils.ParamError(c, "商品ID格式错误")
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 获取原商品
	product, err := productService.GetProduct(productID)
	if err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	// 更新字段
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.OriginalPrice != nil {
		product.OriginalPrice = req.OriginalPrice
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.VideoURL != nil {
		product.VideoURL = *req.VideoURL
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.IsHot != nil {
		product.IsHot = *req.IsHot
	}
	if req.IsNew != nil {
		product.IsNew = *req.IsNew
	}
	if req.Sort != nil {
		product.Sort = *req.Sort
	}

	// 更新图片
	if req.Images != nil {
		product.SetImages(req.Images)
	}

	// 调用服务层
	if err := productService.UpdateProduct(product); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Updated(c, product)
}

// DeleteProduct 删除商品（管理员）
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "商品ID不能为空")
		return
	}

	var productID uint64
	if _, err := fmt.Sscanf(id, "%d", &productID); err != nil {
		utils.ParamError(c, "商品ID格式错误")
		return
	}

	if err := productService.DeleteProduct(productID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, map[string]string{
		"message": "删除成功",
	})
}

// UpdateProductStatus 更新商品状态（管理员）
func UpdateProductStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ParamError(c, "商品ID不能为空")
		return
	}

	var productID uint64
	if _, err := fmt.Sscanf(id, "%d", &productID); err != nil {
		utils.ParamError(c, "商品ID格式错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	product, err := productService.GetProduct(productID)
	if err != nil {
		utils.NotFound(c, "商品不存在")
		return
	}

	product.Status = req.Status

	if err := productService.UpdateProduct(product); err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, map[string]string{
		"message": "状态更新成功",
	})
}
