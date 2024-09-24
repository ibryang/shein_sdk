package api

// AuthParam 基础参数结构体
type AuthParam struct {
	OpenKeyId     string `json:"openKeyId"`
	OpenSecretKey string `json:"secretKey"`
	AppId         string `json:"appid"`
	AppSecretKey  string `json:"appSecretKey"`
}
