package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"hzycoder.com/go-gin-template/internal/auth"
	"hzycoder.com/go-gin-template/internal/config"
	"hzycoder.com/go-gin-template/internal/model"
	"hzycoder.com/go-gin-template/pkg/response"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			response.Abort(ctx, "missing token")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			response.Abort(ctx, "invalid token")
			return
		}

		token := parts[1]

		claims, err := auth.ParseToken([]byte(config.Global.Jwt.Secret), token)
		if err != nil {
			response.Abort(ctx, "invalid token")
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func RequireRole(allowedRoles ...model.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentRole, exists := ctx.Get("role")
		if !exists {
			response.Abort(ctx, "invalid role")
			return
		}

		for _, allowed := range allowedRoles {
			if currentRole == allowed {
				ctx.Next()
				return
			}
		}

		response.Abort(ctx, "invalid role")
	}
}
