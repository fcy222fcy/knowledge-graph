package dto

type PageRequest struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=100"`
}

func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Size
}
