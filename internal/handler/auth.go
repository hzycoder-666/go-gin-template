package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"hzycoder.com/go-gin-template/internal/model"
	dto "hzycoder.com/go-gin-template/internal/model/dto/request"
	"hzycoder.com/go-gin-template/internal/service"
	"hzycoder.com/go-gin-template/pkg/response"
)

func Login(c *gin.Context) {
	var req dto.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		// 存储验证错误到上下文，让中间件处理
		c.Set("validation_error", err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	token, err := service.Login(ctx, req)
	if err != nil {
		slog.Error("login failed", "username", req.Username, "error", err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": token,
	})
}

func Register(c *gin.Context) {
	var req dto.RegisterUser

	if err := c.ShouldBindJSON(&req); err != nil {
		// 存储验证错误到上下文，让中间件处理
		c.Set("validation_error", err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	defaultRole := model.RoleMember
	if req.Role == nil || !model.IsValid(*req.Role) {
		req.Role = &defaultRole
	}

	token, err := service.Register(ctx, req)
	if err != nil {
		slog.Error("register failed", "username", req.Username, "error", err)
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": token,
	})
}
