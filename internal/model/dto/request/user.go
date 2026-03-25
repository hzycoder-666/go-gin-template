package request

type PostUser struct {
	Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Password        string `json:"password" binding:"required,min=8,max=20,password_complexity"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
