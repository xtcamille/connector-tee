//go:build !linux

package main

import (
	"crypto/tls"
)

// GetTLSConfig 返回适用于 Windows/非 Linux 平台的 TLS 配置
// 在这些平台上，我们不进行 enclave 远程度量验证，而是由于连接到的是模拟模式服务器
// 或者由于当前环境不支持 EGo 运行时，所以降级到普通 TLS。
func GetTLSConfig() *tls.Config {
	return &tls.Config{
		// 在本地测试或模拟模式下，可能需要跳过验证
		InsecureSkipVerify: true,
	}
}
