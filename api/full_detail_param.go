package api

// FullDetailParam 商品详情参数结构体
type FullDetailParam struct {
	SkuCodes []string `json:"skuCodes"`
}

type ListParam struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

func (p *ListParam) Default() {
	if p.PageNum == 0 {
		p.PageNum = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 200
	}
}
