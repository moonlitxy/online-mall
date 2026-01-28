package controller

import (
	"online-mall/internal/models"
	"online-mall/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BaseController 基础控制器
type BaseController struct {
}

// GetUserID 获取当前用户ID
func (ctrl *BaseController) GetUserID(c *gin.Context) uint64 {
	userID, _ := c.Get("user_id")
	return userID.(uint64)
}

// GetPagination 获取分页参数
func (ctrl *BaseController) GetPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}

// GetOffset 计算偏移量
func (ctrl *BaseController) GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

// Paginate 分页查询
func (ctrl *BaseController) Paginate(query interface{}, page, pageSize int, total int64) *utils.PageResult {
	return &utils.PageResult{
		List:     query,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

// HandleError 处理错误
func (ctrl *BaseController) HandleError(c *gin.Context, err error, message string) {
	if err != nil {
		utils.ServerError(c)
		return
	}
}

// Success 成功响应
func (ctrl *BaseController) Success(c *gin.Context, data interface{}) {
	utils.Success(c, data)
}

// Error 错误响应
func (ctrl *BaseController) Error(c *gin.Context, code int, message string) {
	utils.Error(c, code, message)
}

// ParamError 参数错误
func (ctrl *BaseController) ParamError(c *gin.Context, message string) {
	utils.ParamError(c, message)
}

// Unauthorized 未授权
func (ctrl *BaseController) Unauthorized(c *gin.Context) {
	utils.Unauthorized(c)
}

// Forbidden 禁止访问
func (ctrl *BaseController) Forbidden(c *gin.Context) {
	utils.Forbidden(c)
}

// NotFound 资源不存在
func (ctrl *BaseController) NotFound(c *gin.Context, message string) {
	utils.NotFound(c, message)
}

// Created 创建成功
func (ctrl *BaseController) Created(c *gin.Context, data interface{}) {
	utils.Created(c, data)
}

// Updated 更新成功
func (ctrl *BaseController) Updated(c *gin.Context, data interface{}) {
	utils.Updated(c, data)
}

// Deleted 删除成功
func (ctrl *BaseController) Deleted(c *gin.Context) {
	utils.Deleted(c)
}

// ValidateError 验证错误
func (ctrl *BaseController) ValidateError(c *gin.Context, errors []string) {
	utils.ValidateError(c, errors)
}

// User 用户控制器
type UserController struct {
	BaseController
	userRepo *models.User
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{}
}

// Product 商品控制器
type ProductController struct {
	BaseController
	productRepo *models.Product
}

// NewProductController 创建商品控制器
func NewProductController() *ProductController {
	return &ProductController{}
}

// Order 订单控制器
type OrderController struct {
	BaseController
	orderRepo *models.Order
}

// NewOrderController 创建订单控制器
func NewOrderController() *OrderController {
	return &OrderController{}
}

// Cart 购物车控制器
type CartController struct {
	BaseController
	cartRepo *models.CartItem
}

// NewCartController 创建购物车控制器
func NewCartController() *CartController {
	return &CartController{}
}

// Address 地址控制器
type AddressController struct {
	BaseController
	addressRepo *models.Address
}

// NewAddressController 创建地址控制器
func NewAddressController() *AddressController {
	return &AddressController{}
}

// Category 分类控制器
type CategoryController struct {
	BaseController
	categoryRepo *models.Category
}

// NewCategoryController 创建分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

// Coupon 优惠券控制器
type CouponController struct {
	BaseController
	couponRepo *models.Coupon
}

// NewCouponController 创建优惠券控制器
func NewCouponController() *CouponController {
	return &CouponController{}
}

// File 文件控制器
type FileController struct {
	BaseController
}

// NewFileController 创建文件控制器
func NewFileController() *FileController {
	return &FileController{}
}

// Search 搜索控制器
type SearchController struct {
	BaseController
}

// NewSearchController 创建搜索控制器
func NewSearchController() *SearchController {
	return &SearchController{}
}
