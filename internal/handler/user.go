package handler

import (
	"strconv"

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
		response.HandleError(c, err)
		return
	}

	response.Success(c, &dto.QueryUser{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

func GetUserInfo(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	userID, ok := ExtractInt64(userIDAny)
	if !ok {
		response.FailWithCode(c, response.CodeParamInvalid)
		return
	}

	ctx := c.Request.Context()
	user, err := service.GetUserById(ctx, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, &dto.QueryUser{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
	})
}

func ExtractInt64(v any) (int64, bool) {
	switch x := v.(type) {
	case int64:
		return x, true
	case int:
		return int64(x), true
	case string:
		if i, err := strconv.ParseInt(x, 10, 64); err == nil {
			return i, true
		}
	case float64:
		return int64(x), true
	}
	return 0, false
}
