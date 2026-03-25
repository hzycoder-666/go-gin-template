package handler

import (
	"github.com/gin-gonic/gin"
	dto "hzycoder.com/go-gin-template/internal/model/dto/request"
	"hzycoder.com/go-gin-template/internal/service"
	"hzycoder.com/go-gin-template/pkg/response"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	ctx := c.Request.Context()

	token, err := service.Login(
		ctx,
		req.Username,
		req.Password,
	)

	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
	})
}

func Register(c *gin.Context) {
	var req dto.PostUser

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	ctx := c.Request.Context()

	token, err := service.Register(ctx, req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
	})
}
