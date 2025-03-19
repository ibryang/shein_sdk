package client

import (
	"context"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/ibryang/shein_sdk/api"
	"github.com/ibryang/shein_sdk/util"
)

// 域名
var Domain = "https://openapi.sheincorp.com"
var DomainTest = "https://openapi-test01.sheincorp.cn"

// buildType 构建类型
type BuildType string

const (
	BuildApp  BuildType = "app"
	BuildOpen BuildType = "open"
)

// SheinClient 客户端结构体
type SheinClient struct {
	Domain     string
	AuthParam  api.AuthParam
	HTTPClient *gclient.Client
}

// NewClient 创建新的 SheinClient 实例
func NewClient() *SheinClient {
	httpClient := gclient.New()
	return &SheinClient{
		Domain:     Domain,
		HTTPClient: httpClient,
	}
}

// SetDomain 设置域名
func (c *SheinClient) SetDomain(domain string) {
	c.Domain = domain
}

// SetAuth 设置认证
func (c *SheinClient) SetAuth(authParam api.AuthParam) {
	c.AuthParam = authParam
}

// PostByAppSign 使用 App 签名进行 POST 请求
func (c *SheinClient) PostByAppSign(path string, data any) (string, error) {
	headers, err := c.getAndBuildSign(BuildApp, path)
	if err != nil {
		return "", err
	}
	return c.post(path, data, headers)
}

// getAndBuildSign 获取时间戳并构建 App 签名所需要的头部
func (c *SheinClient) getAndBuildSign(buildType BuildType, path string) (map[string]string, error) {
	timestamp := util.CurrentTimeMillis()
	secretId := c.AuthParam.OpenKeyId
	secretKey := c.AuthParam.OpenSecretKey
	if buildType == BuildApp {
		secretId = c.AuthParam.AppId
		secretKey = c.AuthParam.AppSecretKey
	}
	sign, err := util.Sign(secretId, secretKey, path, timestamp)
	if err != nil {
		return nil, err
	}
	var headers = map[string]string{}
	if buildType == BuildApp {
		headers["x-lt-appid"] = c.AuthParam.AppId
	} else {
		headers["x-lt-openKeyId"] = c.AuthParam.OpenKeyId
	}
	headers["x-lt-timestamp"] = timestamp
	headers["x-lt-signature"] = sign
	return headers, nil
}

// Post 使用 OpenKey 签名进行 POST 请求
func (c *SheinClient) Post(path string, data any) (string, error) {
	headers, err := c.getAndBuildSign(BuildOpen, path)
	if err != nil {
		return "", err
	}
	return c.post(path, data, headers)
}

// Get 发起 GET 请求
func (c *SheinClient) Get(path string, params any) (string, error) {
	headers, err := c.getAndBuildSign(BuildOpen, path)
	if err != nil {
		return "", err
	}
	return c.get(path, params, headers)
}

// createRequest 创建并配置请求客户端
func (c *SheinClient) createRequest(headers map[string]string) *gclient.Client {
	client := c.HTTPClient.Clone()
	if len(headers) > 0 {
		client.SetHeaderMap(headers)
	}
	return client
}

// handleResponse 处理响应结果
func (c *SheinClient) handleResponse(resp *gclient.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}
	defer resp.Close()
	return resp.ReadAllString(), nil
}

// get 发送 GET 请求
func (c *SheinClient) get(path string, params any, headers map[string]string) (string, error) {
	reqUrl := c.Domain + path

	// 创建请求客户端
	client := c.createRequest(headers)

	// 设置查询参数
	if params != nil {
		mapParams, err := util.ToMapStrStr(params)
		if err != nil {
			return "", err
		}
		for k, v := range mapParams {
			reqUrl += (map[bool]string{true: "?", false: "&"}[!strings.Contains(reqUrl, "?")] + k + "=" + url.QueryEscape(v))
		}
	}

	// 发送请求并获取响应
	resp, err := client.Get(context.Background(), reqUrl)
	return c.handleResponse(resp, err)
}

// post 发送 POST 请求
func (c *SheinClient) post(path string, data any, headers map[string]string) (string, error) {
	reqUrl := c.Domain + path

	// 创建请求客户端
	client := c.createRequest(headers)
	client.SetHeader("Content-Type", "application/json")

	// 发送请求并获取响应（带重试机制）
	resp, err := client.Post(context.Background(), reqUrl, data)
	return c.handleResponse(resp, err)
}
