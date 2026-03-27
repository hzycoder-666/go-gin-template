package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 业务错误码定义
const (
	CodeSuccess      = 0    // 成功
	CodeParamInvalid = 1001 // 参数错误
	CodeUnauthorized = 1002 // 未授权
	CodeForbidden    = 1003 // 禁止访问
	CodeNotFound     = 1004 // 资源不存在

	CodeUserNotFound  = 2001 // 用户不存在
	CodeUserExists    = 2002 // 用户已存在
	CodePasswordWrong = 2003 // 密码错误
	CodeTokenInvalid  = 2004 // Token无效
	CodeTokenExpired  = 2005 // Token过期

	CodeSystemError = 5000 // 系统错误
	CodeDBError     = 5001 // 数据库错误
	CodeCacheError  = 5002 // 缓存错误
)

// 错误信息映射
var errorMessages = map[int]string{
	CodeSuccess:       "success",
	CodeParamInvalid:  "参数错误",
	CodeUnauthorized:  "未授权访问",
	CodeForbidden:     "禁止访问",
	CodeNotFound:      "资源不存在",
	CodeUserNotFound:  "用户不存在",
	CodeUserExists:    "用户已存在",
	CodePasswordWrong: "密码错误",
	CodeTokenInvalid:  "Token无效",
	CodeTokenExpired:  "Token已过期",
	CodeSystemError:   "系统内部错误",
	CodeDBError:       "数据库错误",
	CodeCacheError:    "缓存服务错误",
}

// Resp 统一响应结构
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// BizError 业务错误
type BizError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *BizError) Error() string {
	return e.Message
}

// NewBizError 创建业务错误
func NewBizError(code int, message ...string) *BizError {
	msg := errorMessages[code]
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	return &BizError{
		Code:    code,
		Message: msg,
	}
}

// NewBizErrorWithDetails 创建带详细信息的业务错误
func NewBizErrorWithDetails(code int, details string, message ...string) *BizError {
	err := NewBizError(code, message...)
	err.Details = details
	return err
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code: CodeSuccess,
		Msg:  errorMessages[CodeSuccess],
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Resp{
		Code: CodeParamInvalid,
		Msg:  msg,
	})
}

// FailWithCode 指定错误码的失败响应
func FailWithCode(c *gin.Context, code int, msg ...string) {
	message := errorMessages[code]
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	c.JSON(http.StatusOK, Resp{
		Code: code,
		Msg:  message,
	})
}

// Abort 中止请求（用于认证失败等）
func Abort(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, Resp{
		Code: CodeUnauthorized,
		Msg:  msg,
	})
}

// AbortWithCode 指定状态码的中止请求
func AbortWithCode(c *gin.Context, statusCode, bizCode int, msg ...string) {
	message := errorMessages[bizCode]
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	c.AbortWithStatusJSON(statusCode, Resp{
		Code: bizCode,
		Msg:  message,
	})
}

// HandleError 统一错误处理
func HandleError(c *gin.Context, err error) {
	if bizErr, ok := err.(*BizError); ok {
		// 业务错误
		statusCode := mapBizCodeToHTTPStatus(bizErr.Code)
		c.AbortWithStatusJSON(statusCode, Resp{
			Code: bizErr.Code,
			Msg:  bizErr.Message,
		})
		return
	}

	// 系统错误
	c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{
		Code: CodeSystemError,
		Msg:  errorMessages[CodeSystemError],
	})
}

// mapBizCodeToHTTPStatus 业务错误码映射到HTTP状态码
func mapBizCodeToHTTPStatus(code int) int {
	switch code / 1000 {
	case 1: // 1xxx - 客户端错误
		switch code {
		case CodeUnauthorized:
			return http.StatusUnauthorized
		case CodeForbidden:
			return http.StatusForbidden
		case CodeNotFound:
			return http.StatusNotFound
		default:
			return http.StatusBadRequest
		}
	case 2: // 2xxx - 业务逻辑错误
		return http.StatusBadRequest
	case 5: // 5xxx - 服务端错误
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ParseValidationError 解析验证错误，返回用户友好的错误信息
func ParseValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()

			// 字段名映射为中文
			fieldNameCN := mapFieldName(fieldName)

			// 根据验证标签生成友好的错误信息
			message := generateValidationMessage(fieldNameCN, tag, fieldError.Param())
			messages = append(messages, message)
		}
		return strings.Join(messages, "; ")
	}

	// 如果不是验证错误，返回原始错误
	return err.Error()
}

// mapFieldName 字段名映射为中文
func mapFieldName(fieldName string) string {
	fieldMap := map[string]string{
		"Username":        "用户名",
		"Password":        "密码",
		"ConfirmPassword": "确认密码",
		"Email":           "邮箱",
		"Phone":           "手机号",
		"Role":            "角色",
		"Nickname":        "昵称",
	}
	if cn, exists := fieldMap[fieldName]; exists {
		return cn
	}
	return fieldName
}

// generateValidationMessage 生成验证错误信息
func generateValidationMessage(fieldName, tag, param string) string {
	switch tag {
	case "required":
		return fieldName + "不能为空"
	case "min":
		return fieldName + "长度不能少于" + param + "个字符"
	case "max":
		return fieldName + "长度不能超过" + param + "个字符"
	case "email":
		return fieldName + "格式不正确"
	case "eqfield":
		return fieldName + "与" + mapFieldName(param) + "不一致"
	case "alphanum":
		return fieldName + "只能包含字母和数字"
	case "password_complexity":
		return fieldName + "必须包含大写字母、小写字母、数字和特殊字符"
	default:
		return fieldName + "格式不正确"
	}
}
