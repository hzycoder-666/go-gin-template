package response

import "hzycoder.com/go-gin-template/internal/model"

type QueryUser struct {
	ID       int64      `json:"id"`
	Username string     `json:"username"`
	Nickname *string    `json:"nickname"`
	Role     model.Role `json:"role"`
}
