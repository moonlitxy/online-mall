package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"online-mall/internal/config"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint64, username string, role string) (string, error) {
	cfg := config.GlobalConfig.JWT

	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)

	// 创建声明
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
			Subject:   username,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*JWTClaims, error) {
	cfg := config.GlobalConfig.JWT

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// 提取声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新JWT token
func RefreshToken(tokenString string) (string, error) {
	// 解析token获取声明
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新token
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}

// ValidateToken 验证JWT token
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// 直接使用ParseToken进行验证
	return ParseToken(tokenString)
}

// GetUserIDFromToken 从token中获取用户ID
func GetUserIDFromToken(tokenString string) (uint64, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetUsernameFromToken 从token中获取用户名
func GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// IsTokenExpired 检查token是否过期
func IsTokenExpired(tokenString string) bool {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return true
	}
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetTokenExpiration 获取token过期时间
func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}

// CreateTokenWithContext 在context中创建token（用于中间件）
func CreateTokenWithContext(userID uint64, username string, role string) (string, error) {
	return GenerateToken(userID, username, role)
}
