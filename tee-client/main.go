package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/xtcamille/connector-tee/api"
)

func main() {
	// 1. 建立与服务器的 TLS 连接

	// 在 Occlum 环境下，如果是 RA-TLS，通常需要验证服务器的证明。
	// 这里我们使用标准 TLS 配置。如果需要跳过验证（仅用于开发）或提供根 CA。

	insecure := os.Getenv("INSECURE") == "1"
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure,
	}

	// 启动地址
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "localhost:9001"
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to dial %s: %v", addr, err)
	}
	defer conn.Close()

	if insecure {
		fmt.Println("Connection established (Insecure - verification skipped).")
	} else {
		fmt.Println("Connection established (TLS verified).")
	}

	// 2. 构造请求
	code := `
        function main(params) {
            let a = params.a;
            let b = params.b;
            return { sum: a + b, product: a * b };
        }
    `
	params := `{"a": 10, "b": 20}`

	req := api.ExecuteRequest{
		Code:   code,
		Params: params,
	}

	// 3. 发送请求
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		log.Fatalf("Failed to encode request: %v", err)
	}

	// 4. 读取响应
	var resp api.ExecuteResponse
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&resp); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Error != "" {
		log.Fatalf("Execution returned error: %s", resp.Error)
	}

	fmt.Printf("Result: %s\n", resp.Result)
}
