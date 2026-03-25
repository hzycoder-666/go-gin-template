package response

import "github.com/gin-gonic/gin"

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(200, Resp{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(200, Resp{
		Code: 1,
		Msg:  msg,
	})
}
