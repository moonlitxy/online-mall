package utils

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data: PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context, message string) {
	Error(c, 400, message)
}

// BadRequest 错误请求
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	Error(c, 401, "未授权")
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context) {
	Error(c, 403, "禁止访问")
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// ServerError 服务器错误
func ServerError(c *gin.Context) {
	Error(c, 500, "服务器内部错误")
}

// TooManyRequests 请求过于频繁
func TooManyRequests(c *gin.Context) {
	Error(c, 429, "请求过于频繁，请稍后再试")
}

// ValidateError 验证错误
func ValidateError(c *gin.Context, errors []string) {
	c.JSON(http.StatusOK, Response{
		Code:    400,
		Message: "请求参数验证失败",
		Data: gin.H{
			"errors": errors,
		},
	})
}

// Created 创建成功
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "created successfully",
		Data:    data,
	})
}

// Updated 更新成功
func Updated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "updated successfully",
		Data:    data,
	})
}

// Deleted 删除成功
func Deleted(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "deleted successfully",
		Data:    nil,
	})
}

// UploadSuccess 上传成功
func UploadSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "upload successfully",
		Data:    data,
	})
}

// DownloadSuccess 下载成功
func DownloadSuccess(c *gin.Context, data interface{}, filename string) {
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data.([]byte))
}

// StreamSuccess 流式响应
func StreamSuccess(c *gin.Context, reader io.Reader) {
	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return false
		}
		if n > 0 {
			w.Write(buf[:n])
		}
		return err == nil
	})
}

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	Type string      `json:"type"` // 消息类型
	Data interface{} `json:"data"` // 消息数据
	From string      `json:"from"` // 发送者
	To   string      `json:"to"`   // 接收者
	Time int64       `json:"time"` // 时间戳
}

// SendWebSocketMessage 发送WebSocket消息
func SendWebSocketMessage(c *gin.Context, message WebSocketMessage) {
	c.JSON(http.StatusOK, message)
}
