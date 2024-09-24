package client

import (
	"github.com/ibryang/shein_sdk/api"
)

// GoodsClient 商品客户端
type GoodsClient struct {
	*SheinClient
	fullDetailPath string
	listPath       string
}

// NewGoodsClient 创建新的 GoodsClient 实例
func NewGoodsClient() *GoodsClient {
	return &GoodsClient{
		SheinClient:    NewClient(),
		fullDetailPath: "/open-api/openapi-business-backend/product/full-detail",
		listPath:       "/open-api/openapi-business-backend/product/query",
	}
}

// FullDetail 获取商品详情
func (c *GoodsClient) FullDetail(param *api.FullDetailParam) (string, error) {
	return c.Post(c.fullDetailPath, param)
}

// List 获取商品列表
func (c *GoodsClient) List(param *api.ListParam) (string, error) {
	if param == nil {
		param = &api.ListParam{}
		param.Default()
	}
	return c.Post(c.listPath, param)
}
