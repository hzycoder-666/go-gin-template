package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"hzycoder.com/go-gin-template/pkg/response"
)

// RecoveryWithBizError 兼容 gin.Recovery 但使用业务错误格式
func RecoveryWithBizError() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		slog.Error("panic recovered",
			"error", fmt.Sprintf("%+v", recovered),
			"stack", string(debug.Stack()))
		response.AbortWithCode(c, 500, response.CodeSystemError, "服务器内部错误")
	})
}

// BizErrorHandler 统一业务错误和验证错误处理
func BizErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 已经写了响应就无需二次处理
		if c.Writer.Written() {
			return
		}

		// 验证错误优先
		if errVal, exists := c.Get("validation_error"); exists {
			if err, ok := errVal.(error); ok {
				msg := response.ParseValidationError(err)
				response.FailWithCode(c, response.CodeParamInvalid, msg)
				return
			}
		}

		// Gin errors 处理（从 service/handler 传过来的 BizError）
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			response.HandleError(c, err)
			return
		}
	}
}
