package utils

import (
	"online-mall/internal/models"
)

// InitDB 初始化数据库
func InitDB() error {
	return models.InitDatabase()
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	return models.CloseDatabase()
}
