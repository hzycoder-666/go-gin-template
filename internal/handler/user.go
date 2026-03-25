package handler

import (
	dto "hzycoder.com/go-gin-template/internal/model/dto/response"
	"hzycoder.com/go-gin-template/internal/service"
	"hzycoder.com/go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	username := c.Param("username")

	ctx := c.Request.Context()
	user, err := service.GetUser(ctx, username)

	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, &dto.QueryUser{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Nickname: user.Nickname,
	})
}

func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	response.Success(c, gin.H{
		"user_id": userID,
	})
}
