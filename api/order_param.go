package api

import (
	"time"
)

// 采购单列表参数
type PurchaseOrderInfosParam struct {
	PageNumber        int    `json:"pageNumber"`                  // 页码
	PageSize          int    `json:"pageSize"`                    // 一页大小，最大200条
	OrderNos          string `json:"orderNos,omitempty"`          // 采购单号，支持批量，一次最多请求200采购单号
	Skcs              string `json:"skcs,omitempty"`              // skc列表
	Type              string `json:"type,omitempty"`              // 订单类型数组（1急采，2备货）
	SupplierCodes     string `json:"supplierCodes,omitempty"`     // 商家货号数组
	AllocateTimeStart string `json:"allocateTimeStart,omitempty"` // 查询备货单开始时间;分单时间-开始(Date格式) 如：2018-05-23 10:29:59； 查询不能超过60天
	AllocateTimeEnd   string `json:"allocateTimeEnd,omitempty"`   // 查询备货单结束时间;分单时间-结束(Date格式),如：2018-06-23 10:29:59 ；查询不能超过60天
	AddTimeStart      string `json:"addTimeStart,omitempty"`      // 查询急采单开始时间;下单时间-开始(Date格式), 如：2018-05-23 10:29:59； 查询不能超过60天
	AddTimeEnd        string `json:"addTimeEnd,omitempty"`        // 查询急采单结束时间;下单时间-结束(Date格式),如2018-06-23 10:29:59； 查询不能超过60天
	CombineTimeStart  string `json:"combineTimeStart,omitempty"`  // 查询急采单和备货单开始时间;联合时间-开始(Date格式), 如：2018-05-23 10:29:59 ；查询不能超过60天
	CombineTimeEnd    string `json:"combineTimeEnd,omitempty"`    // 查询急采单和备货单结束时间;联合时间-结束(Date格式),如2018-06-23 10:29:59 ；查询不能超过60天
	UpdateTimeStart   string `json:"updateTimeStart,omitempty"`   // 查询状态发生变更的采购单开始时间；更新时间-开始(Date格式),如2018-06-23 10:29:59 ；查询不能超过60天
	UpdateTimeEnd     string `json:"updateTimeEnd,omitempty"`     // 查询状态发生变更的采购单结束时间；更新时间-结束(Date格式),如2018-06-23 10:29:59 ；查询不能超过60天
	SelectJitMother   int    `json:"selectJitMother,omitempty"`   // 是否查询jit母单，1-是--返回JIT母单和非JIT母单信息；2或者不填--否-只返回非JIT母单信息；简易平台无需关注
}

func (p *PurchaseOrderInfosParam) Default() {
	p.PageNumber = 1
	p.PageSize = 200
	p.Type = "1"
	p.UpdateTimeStart = time.Now().AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	p.UpdateTimeEnd = time.Now().Format("2006-01-02 15:04:05")
}

// 设置更新时间 参数: hours 小时 可以是浮点数
func (p *PurchaseOrderInfosParam) SetUpdateTime(hours float64) {
	// 将小数转为分钟
	minutes := hours * 60
	p.UpdateTimeStart = time.Now().Add(-time.Duration(minutes) * time.Minute).Format("2006-01-02 15:04:05")
	p.UpdateTimeEnd = time.Now().Format("2006-01-02 15:04:05")
}
