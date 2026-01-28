package models

import "errors"

var (
	// ErrRecordNotFound 记录未找到
	ErrRecordNotFound = errors.New("record not found")

	// ErrDuplicateKey 重复键错误
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrForeignKeyConstraint 外键约束错误
	ErrForeignKeyConstraint = errors.New("foreign key constraint")

	// ErrInvalidArgument 无效参数错误
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = errors.New("permission denied")

	// ErrTokenExpired Token过期
	ErrTokenExpired = errors.New("token expired")

	// ErrTokenInvalid 无效Token
	ErrTokenInvalid = errors.New("invalid token")
)

// IsRecordNotFound 检查是否是记录未找到错误
func IsRecordNotFound(err error) bool {
	return errors.Is(err, ErrRecordNotFound)
}

// IsDuplicateKey 检查是否是重复键错误
func IsDuplicateKey(err error) bool {
	return errors.Is(err, ErrDuplicateKey)
}
