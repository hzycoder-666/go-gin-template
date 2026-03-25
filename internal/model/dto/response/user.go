package response

type QueryUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
}
