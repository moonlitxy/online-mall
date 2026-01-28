package controller

import (
	"online-mall/internal/models"
	"online-mall/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"omitempty,e164"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 查找用户
	user := &models.User{}
	if err := models.DB.Where("username = ? OR phone = ? OR email = ?", req.Username, req.Username, req.Username).First(user).Error; err != nil {
		if err == models.ErrRecordNotFound {
			utils.ParamError(c, "用户名或密码错误")
			return
		}
		utils.ServerError(c)
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		utils.ParamError(c, "账号已被禁用")
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		utils.ParamError(c, "用户名或密码错误")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.ServerError(c)
		return
	}

	// 更新最后登录时间
	now := time.Now()
	models.DB.Model(user).Update("last_login_at", now)

	// 获取用户信息
	userInfo := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"phone":    user.Phone,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"role":     "user",
	}

	utils.Success(c, map[string]interface{}{
		"token": token,
		"user":  userInfo,
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 检查用户名是否已存在
	var count int64
	if err := models.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		utils.ServerError(c)
		return
	}
	if count > 0 {
		utils.ParamError(c, "用户名已存在")
		return
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		if err := models.DB.Model(&models.User{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
			utils.ServerError(c)
			return
		}
		if count > 0 {
			utils.ParamError(c, "手机号已被注册")
			return
		}
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if err := models.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
			utils.ServerError(c)
			return
		}
		if count > 0 {
			utils.ParamError(c, "邮箱已被注册")
			return
		}
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Phone:    req.Phone,
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}

	if err := models.DB.Create(user).Error; err != nil {
		utils.ServerError(c)
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.ServerError(c)
		return
	}

	// 返回用户信息（不包含密码）
	userInfo := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"phone":    user.Phone,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"role":     "user",
	}

	utils.Created(c, map[string]interface{}{
		"token": token,
		"user":  userInfo,
	})
}

// RefreshToken 刷新token
func RefreshToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		utils.Unauthorized(c)
		return
	}

	// 解析token
	claims, err := utils.ParseToken(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		utils.Unauthorized(c)
		return
	}

	// 生成新token
	newToken, err := utils.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, map[string]string{
		"token": newToken,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 获取用户ID
	_, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c)
		return
	}

	// 这里可以将token加入黑名单，或者在Redis中标记token为无效
	// 为了简化，这里只返回成功响应
	utils.Success(c, nil)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	user := &models.User{}
	if err := models.DB.First(user, userID).Error; err != nil {
		if err == models.ErrRecordNotFound {
			utils.NotFound(c, "用户不存在")
			return
		}
		utils.ServerError(c)
		return
	}

	// 返回用户信息（不包含密码）
	userInfo := map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"nickname":      user.Nickname,
		"phone":         user.Phone,
		"email":         user.Email,
		"avatar":        user.Avatar,
		"status":        user.Status,
		"last_login_at": user.LastLoginAt,
		"created_at":    user.CreatedAt,
	}

	utils.Success(c, userInfo)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	var req struct {
		Nickname string `json:"nickname" binding:"max=50"`
		Phone    string `json:"phone" binding:"omitempty,e164"`
		Email    string `json:"email" binding:"omitempty,email"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 检查手机号是否已被其他用户使用
	if req.Phone != "" {
		var count int64
		if err := models.DB.Model(&models.User{}).
			Where("phone = ? AND id != ?", req.Phone, userID).
			Count(&count).Error; err != nil {
			utils.ServerError(c)
			return
		}
		if count > 0 {
			utils.ParamError(c, "手机号已被使用")
			return
		}
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" {
		var count int64
		if err := models.DB.Model(&models.User{}).
			Where("email = ? AND id != ?", req.Email, userID).
			Count(&count).Error; err != nil {
			utils.ServerError(c)
			return
		}
		if count > 0 {
			utils.ParamError(c, "邮箱已被使用")
			return
		}
	}

	// 更新用户信息
	user := &models.User{}
	user.ID = uint64(userID)
	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) > 0 {
		if err := models.DB.Model(user).Updates(updates).Error; err != nil {
			utils.ServerError(c)
			return
		}
	}

	// 重新获取用户信息
	if err := models.DB.First(user, userID).Error; err != nil {
		utils.ServerError(c)
		return
	}

	// 返回更新后的用户信息
	userInfo := map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"nickname":      user.Nickname,
		"phone":         user.Phone,
		"email":         user.Email,
		"avatar":        user.Avatar,
		"status":        user.Status,
		"last_login_at": user.LastLoginAt,
		"updated_at":    user.UpdatedAt,
	}

	utils.Updated(c, userInfo)
}

// UpdatePassword 更新密码
func UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.Unauthorized(c)
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required,min=6"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, "请求参数格式错误")
		return
	}

	// 获取用户信息
	user := &models.User{}
	if err := models.DB.First(user, userID).Error; err != nil {
		utils.ServerError(c)
		return
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		utils.ParamError(c, "旧密码错误")
		return
	}

	// 更新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ServerError(c)
		return
	}

	if err := models.DB.Model(user).Update("password", string(hashedPassword)).Error; err != nil {
		utils.ServerError(c)
		return
	}

	utils.Success(c, map[string]string{
		"message": "密码更新成功",
	})
}
