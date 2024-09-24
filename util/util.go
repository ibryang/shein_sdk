package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ToJSONString 将对象转换为 JSON 字符串
func ToJSONString(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// ToMapStrAny 将对象转换为 Map[string]interface{}
func ToMapStrAny(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToMapStrStr(data interface{}) (map[string]string, error) {
	mapData, err := ToMapStrAny(data)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for key, value := range mapData {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result, nil
}

// ReadResponseBody 读取响应体
func ReadResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
