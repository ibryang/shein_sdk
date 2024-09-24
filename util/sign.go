package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// SignUtil 签名工具
const (
	HMAC_SHA256   = "sha256"
	RANDOM_LENGTH = 5
)

// Sign 生成签名
func Sign(keyID, secretKey, apiPath, timestamp string) (string, error) {
	signString := fmt.Sprintf("%s&%s&%s", keyID, timestamp, apiPath)
	randomKey, err := generateRandomKey(RANDOM_LENGTH)
	if err != nil {
		return "", err
	}
	hashValue, err := hmacSha256(signString, secretKey+randomKey)
	if err != nil {
		return "", err
	}
	hashValueString := hex.EncodeToString(hashValue)
	base64HashValueString := base64.StdEncoding.EncodeToString([]byte(hashValueString))
	signature := randomKey + base64HashValueString
	return signature, nil
}

// VerifySign 验证签名
func VerifySign(signature, keyID, secretKey, requestPath, timestamp string) (bool, error) {
	if len(signature) < RANDOM_LENGTH {
		return false, errors.New("signature length is too short")
	}
	randomKey := signature[:RANDOM_LENGTH]
	signString := fmt.Sprintf("%s&%s&%s", keyID, timestamp, requestPath)
	hashValue, err := hmacSha256(signString, secretKey+randomKey)
	if err != nil {
		return false, err
	}
	hashValueString := hex.EncodeToString(hashValue)
	base64HashValueString := base64.StdEncoding.EncodeToString([]byte(hashValueString))
	base64Value := randomKey + base64HashValueString
	return signature == base64Value, nil
}

// hmacSha256 生成 HMAC-SHA256 哈希
func hmacSha256(message, secret string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

// generateRandomKey 生成随机键
func generateRandomKey(length int) (string, error) {
	const str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	data := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	if len(data) < length {
		return "", errors.New("not enough data to generate random key")
	}
	randomKey := ""
	for i := 0; i < length; i++ {
		index := int(data[i]) % len(str)
		randomKey += string(str[index])
	}
	return randomKey, nil
}
