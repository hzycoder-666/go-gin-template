package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		duration := time.Since(start)

		slog.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.Query(),
			"status", c.Writer.Status(),
			"latency", duration,
		)
	}
}
