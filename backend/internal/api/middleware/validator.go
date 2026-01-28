package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans    ut.Translator
	validate *validator.Validate
)

// init 初始化验证器
func init() {
	// 设置中文翻译器
	zhTrans := zh.New()
	trans, _ = ut.New(zhTrans, zhTrans).GetTranslator("zh")

	validate = validator.New()

	// 注册验证器的翻译
	translations.RegisterDefaultTranslations(validate, trans)

	// 注册自定义验证函数
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} 手机号格式不正确", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})
}

// ValidateRequest 请求验证中间件
func ValidateRequest(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定请求体
		if err := c.ShouldBindJSON(model); err != nil {
			// 如果是验证错误，返回详细的中文错误信息
			if ve, ok := err.(validator.ValidationErrors); ok {
				var errorMsgs []string
				for _, e := range ve {
					errorMsgs = append(errorMsgs, e.Translate(trans))
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "请求参数验证失败",
					"errors":  errorMsgs,
				})
				c.Abort()
				return
			}

			// 其他类型的错误
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数解析失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将验证后的数据存入上下文
		c.Set("validated_data", model)

		c.Next()
	}
}

// ValidateQuery 查询参数验证中间件
func ValidateQuery(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定查询参数
		if err := c.ShouldBindQuery(model); err != nil {
			if ve, ok := err.(validator.ValidationErrors); ok {
				var errorMsgs []string
				for _, e := range ve {
					errorMsgs = append(errorMsgs, e.Translate(trans))
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "查询参数验证失败",
					"errors":  errorMsgs,
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "查询参数解析失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("validated_query", model)
		c.Next()
	}
}

// GetValidatedData 从上下文中获取验证后的数据
func GetValidatedData(c *gin.Context, key string) interface{} {
	if data, exists := c.Get(key); exists {
		return data
	}
	return nil
}

// ValidatePhone 手机号验证
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return validateChinesePhone(phone)
}

// validateChinesePhone 验证中国手机号
func validateChinesePhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}

	// 检查是否全是数字
	if phone[0] != '1' {
		return false
	}

	// 检查运营商前缀
	prefixes := []string{
		"130", "131", "132", "133", "134", "135", "136", "137", "138", "139",
		"145", "147", "149",
		"150", "151", "152", "153", "155", "156", "157", "158", "159",
		"165", "166", "167",
		"170", "171", "172", "173", "174", "175", "176", "177", "178",
		"180", "181", "182", "183", "184", "185", "186", "187", "188", "189",
		"191", "198", "199",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(phone, prefix) {
			return true
		}
	}

	return false
}

// ValidateEmail 邮箱验证
func ValidateEmail(email string) bool {
	// 简单的邮箱格式验证
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ValidatePassword 密码验证
func ValidatePassword(password string) bool {
	// 密码长度至少6位，包含字母和数字
	if len(password) < 6 {
		return false
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}

// SanitizeInput 输入清理中间件
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 清理查询参数
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				c.Request.URL.Query().Set(key, strings.TrimSpace(values[0]))
			}
		}

		// 清理表单数据
		if c.Request.Form != nil {
			for key, values := range c.Request.Form {
				if len(values) > 0 {
					c.Request.Form.Set(key, strings.TrimSpace(values[0]))
				}
			}
		}

		c.Next()
	}
}

// RateLimiter 限流中间件（简化版本）
func RateLimiter() gin.HandlerFunc {
	// 这里可以使用Redis实现限流
	// 为了简化，使用内存map
	requests := make(map[string]int)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 简单的计数器实现
		requests[ip]++

		if requests[ip] > 100 { // 每IP每分钟100次请求
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
