package util

import (
	"strconv"
	"time"
)

// CurrentTimeMillis 获取当前时间的毫秒数
func CurrentTimeMillis() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}
