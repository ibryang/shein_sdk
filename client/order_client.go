package client

import (
	"github.com/ibryang/shein_sdk/api"
)

// OrderClient 订单客户端
type OrderClient struct {
	*SheinClient
	purchaseOrderInfosPath string
}

// NewOrderClient 创建新的 OrderClient 实例
func NewOrderClient() *OrderClient {
	return &OrderClient{
		SheinClient:            NewClient(),
		purchaseOrderInfosPath: "/open-api/order/purchase-order-infos",
	}
}

// PurchaseOrderInfos 获取采购单列表
func (c *OrderClient) PurchaseOrderInfos(param *api.PurchaseOrderInfosParam) (string, error) {
	if param == nil {
		param = &api.PurchaseOrderInfosParam{}
		param.Default()
	}
	return c.Get(c.purchaseOrderInfosPath, param)
}
