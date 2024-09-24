package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SheinClient 客户端结构体
type SheinClient struct {
	Domain     string
	AuthParam  api.AuthParam
	HTTPClient HTTPClient
}

// NewClient 创建新的 SheinClient 实例
func NewClient() *SheinClient {
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

// get 发送 GET 请求
func (c *SheinClient) get(path string, params any, headers map[string]string) (string, error) {
	// 拼接参数到 URL
	if params != nil {
		mapParams, err := util.ToMapStrStr(params)
		if err != nil {
			return "", err
		}
		query := ""
		i := 0
		for key, value := range mapParams {
			value = url.QueryEscape(value)
			if i == 0 {
				query += "?" + key + "=" + value
			} else {
				query += "&" + key + "=" + value
			}
			i++
		}
		path += query
	}

	url := c.Domain + path

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 设置头部
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 打印调试信息
	fmt.Printf("GET请求地址: %s\n", path)
	fmt.Printf("请求header: %s\n", util.ToJSONString(headers))

	// 执行请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return "", err
	}

	return body, nil
}

// post 发送 POST 请求
func (c *SheinClient) post(path string, data any, headers map[string]string) (string, error) {
	url := c.Domain + path
	var jsonData []byte
	var err error
	if data != nil {
		jsonData, err = json.Marshal(data)
		if err != nil {
			return "", err
		}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置头部
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// 打印调试信息
	fmt.Printf("POST请求地址: %s\n", url)
	// fmt.Printf("请求header: %s\n", util.ToJSONString(headers))
	fmt.Printf("请求参数: %s\n", string(jsonData))

	// 执行请求
	resp, err := c.doRequestWithRetry(req, 3)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := util.ReadResponseBody(resp)
	if err != nil {
		return "", err
	}

	// 打印响应
	// fmt.Printf("response: %s\n", body)

	return body, nil
}

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
}

// doRequestWithRetry 发送请求并重试
func (c *SheinClient) doRequestWithRetry(req *http.Request, retries int) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < retries; i++ {
		resp, err = httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}
	return resp, err
}
