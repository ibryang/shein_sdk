package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AesUtil AES 加密工具
const (
	DefaultIVSeed = "space-station-de"
	CipherAlgo    = "AES-128-CBC"
)

// Decrypt 解密内容
func Decrypt(content, key string) (string, error) {
	block, err := aes.NewCipher([]byte(content)[:aes.BlockSize])
	if err != nil {
		return "", err
	}
	iv := []byte("space-station-de")[:aes.BlockSize]
	cbcDecrypter := cipher.NewCBCDecrypter(block, iv)
	encryptedContentBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	decryptedContent := make([]byte, len(encryptedContentBytes))
	cbcDecrypter.CryptBlocks(decryptedContent, encryptedContentBytes)
	decryptedContent = PKCS7Unpad(decryptedContent)
	return string(decryptedContent), nil
}

// PKCS7Unpad 去除填充
func PKCS7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding < 1 || unpadding > length {
		return nil
	}
	return data[:(length - unpadding)]
}
