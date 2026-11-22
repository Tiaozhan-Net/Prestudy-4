package common

type IDUriForm struct {
	ID int `uri:"id" json:"id" binding:"min=1"`
}

// PagerForm 通用分页参数：入参绑定 + 响应内嵌（json tag 保证输出 page/size）
type PagerForm struct {
	Page int `form:"page" json:"page" binding:"omitempty,min=1"`
	Size int `form:"size" json:"size" binding:"omitempty,min=1,max=100"`
}

// Normalize 兜底：page 至少 1、size 默认 10、上限 100
func (p PagerForm) Normalize() PagerForm {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return p
}

// Offset 与 model.Paginate 口径一致
func (p PagerForm) Offset() int {
	return (p.Page - 1) * p.Size
}
