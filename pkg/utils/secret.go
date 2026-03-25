package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func MakeSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	// 使用 URL 安全的 Base64 编码（无 + /，适合环境变量）
	secret := base64.URLEncoding.EncodeToString(bytes)
	return secret
}
