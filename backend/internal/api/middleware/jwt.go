package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-mall/internal/utils"
	"strings"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 继续处理请求
		c.Next()
	}
}

// GetCurrentUser 从context中获取当前用户信息
func GetCurrentUser(c *gin.Context) (userID uint64, username string, role string, ok bool) {
	userIDValue, userIDOk := c.Get("user_id")
	usernameValue, usernameOk := c.Get("username")
	roleValue, roleOk := c.Get("role")

	if userIDOk && usernameOk && roleOk {
		userID, _ = userIDValue.(uint64)
		username, _ = usernameValue.(string)
		role, _ = roleValue.(string)
		return userID, username, role, true
	}

	return 0, "", "", false
}

// RequireAdmin 要求管理员权限
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, role, ok := GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选认证（允许未登录用户访问）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有token，继续处理请求
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// token格式错误，继续处理请求
			c.Next()
			return
		}

		tokenString := parts[1]

		// 验证token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			// token无效，继续处理请求
			c.Next()
			return
		}

		// 将用户信息存入context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
