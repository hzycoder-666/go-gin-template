package request

type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`              // URL query 参数名: page
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"` // 推荐用下划线命名
}
