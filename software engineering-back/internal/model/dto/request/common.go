package request

// PageRequest 分页请求参数
type PageRequest struct {
	Page int `form:"page" binding:"required,min=1"`     // 页码，从1开始
	Size int `form:"size" binding:"required,min=1,max=100"` // 每页数量，1-100
}

// Offset 计算数据库查询的偏移量
func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Size
}
