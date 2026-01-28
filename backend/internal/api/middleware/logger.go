package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"online-mall/internal/config"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 请求路径
		path := c.Request.URL.Path

		// 请求方法
		method := c.Request.Method

		// 请求IP
		clientIP := c.ClientIP()

		// 请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()

		// 响应状态码
		statusCode := c.Writer.Status()

		// 响应体大小
		responseSize := c.Writer.Size()

		// 执行时间
		latency := end.Sub(start)

		// 错误信息
		var errors []string
		for _, err := range c.Errors {
			errors = append(errors, err.Error())
		}

		// 构建日志条目
		logEntry := map[string]interface{}{
			"client_ip":     clientIP,
			"method":        method,
			"path":          path,
			"query":         c.Request.URL.RawQuery,
			"status":        statusCode,
			"latency":       latency,
			"latency_human": latency.String(),
			"response_size": responseSize,
			"request_body":  string(requestBody),
			"errors":        errors,
			"start_time":    start.Format("2006-01-02 15:04:05"),
			"end_time":      end.Format("2006-01-02 15:04:05"),
		}

		// 根据状态码记录日志级别
		if statusCode >= 500 {
			log.Printf("[ERROR] %s %s %d %s", method, path, statusCode, latency)
		} else if statusCode >= 400 {
			log.Printf("[WARN] %s %s %d %s", method, path, statusCode, latency)
		} else {
			log.Printf("[INFO] %s %s %d %s", method, path, statusCode, latency)
		}

		// 记录详细日志（开发环境）
		if config.GlobalConfig.App.Debug {
			logJSON, _ := json.Marshal(logEntry)
			log.Printf("[DEBUG] %s", string(logJSON))
		}
	}
}

// LoggerToFile 日志到文件的中间件
func LoggerToFile() gin.HandlerFunc {
	_ = config.GlobalConfig.Log

	// 这里可以集成第三方日志库如 logrus、zap 等
	// 为了简化，这里使用标准库的 log

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()

		// 记录日志到文件
		log.Printf("[%s] %s %s %d %v",
			method,
			path,
			c.ClientIP(),
			statusCode,
			latency,
		)
	}
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := generateRequestID()

		// 设置响应头
		c.Header("X-Request-ID", requestID)

		// 存入上下文
		c.Set("request_id", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetRequestID 从上下文中获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// Recovery 恢复中间件，防止panic
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic错误
				log.Printf("[PANIC] %v", err)

				// 返回错误响应
				c.JSON(500, gin.H{
					"code":       500,
					"message":    "Internal Server Error",
					"request_id": GetRequestID(c),
				})

				// 终止请求
				c.Abort()
			}
		}()

		c.Next()
	}
}
