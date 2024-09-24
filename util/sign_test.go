package util

import (
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	// 示例：生成签名
	sign, err := Sign("1000000", "test123", "/open-api/order/purchase-order-infos", "1583398764000")
	if err != nil {
		fmt.Println("Sign Error:", err)
		return
	}
	fmt.Println("Generated Sign:", sign)
}
