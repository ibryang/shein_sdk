package client

import (
	"github.com/ibryang/shein_sdk/api"
)

// GetTempTokenClient 获取临时 Token 客户端
type GetTempTokenClient struct {
	*SheinClient
	path string
}

// NewGetTempTokenClient 创建新的 GetByTempTokenClient 实例
func NewGetTempTokenClient() *GetTempTokenClient {
	return &GetTempTokenClient{
		SheinClient: NewClient(),
		path:        "/open-api/auth/get-by-token",
	}
}

// GetTempToken 获取临时 Token
func (c *GetTempTokenClient) GetTempToken(param api.GetByTokenParam) (string, error) {
	return c.PostByAppSign(c.path, param)
}
