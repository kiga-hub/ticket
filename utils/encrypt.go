package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// 加密十六进制hash字符串
func HashHmacHexDigest(data, key string) string {
	return fmt.Sprintf("%x\n", HashHmacDigest(data, key))
}

func HashHmacDigest(data, key string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

// 数据base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
