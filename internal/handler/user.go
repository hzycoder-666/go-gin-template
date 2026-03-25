package handler

import (
	"hzycoder.com/go-gin-template/internal/service"
	"hzycoder.com/go-gin-template/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := service.GetUser(username)

	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, user)
}

func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	response.Success(c, gin.H{
		"user_id": userID,
	})
}
