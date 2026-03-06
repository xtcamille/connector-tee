package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/dop251/goja"
	"github.com/xtcamille/connector-tee/api"
)

// Execute 在 enclave 内部执行用户提供的 JavaScript 代码
func Execute(req *api.ExecuteRequest) (*api.ExecuteResponse, error) {
	// 1. 创建 goja 运行时（沙箱）
	vm := goja.New()

	// 2. 注入安全的 console.log
	vm.Set("console", map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			for _, arg := range call.Arguments {
				fmt.Print(arg.String(), " ")
			}
			fmt.Println()
			return nil
		},
	})

	// 3. 编译代码
	_, err := vm.RunScript("main.js", req.Code)
	if err != nil {
		return nil, fmt.Errorf("compile error: %v", err)
	}

	// 4. 获取 main 函数
	var mainFunc func(interface{}) interface{}
	err = vm.ExportTo(vm.Get("main"), &mainFunc)
	if err != nil {
		return nil, fmt.Errorf("main function not found: %v", err)
	}

	// 5. 解析参数
	var params interface{}
	if req.Params != "" {
		if err := json.Unmarshal([]byte(req.Params), &params); err != nil {
			return nil, fmt.Errorf("invalid params JSON: %v", err)
		}
	}

	// 6. 带超时执行
	type result struct {
		value interface{}
		err   error
	}
	done := make(chan result, 1)

	timer := time.AfterFunc(5*time.Second, func() {
		vm.Interrupt("execution timeout")
	})
	defer timer.Stop()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- result{err: fmt.Errorf("panic: %v", r)}
			}
		}()
		res := mainFunc(params)
		done <- result{value: res}
	}()

	select {
	case res := <-done:
		if res.err != nil {
			return &api.ExecuteResponse{Error: fmt.Sprintf("execution error: %v", res.err)}, nil
		}
		// 7. 将结果序列化为 JSON
		resultBytes, err := json.Marshal(res.value)
		if err != nil {
			return &api.ExecuteResponse{Error: fmt.Sprintf("result marshal error: %v", err)}, nil
		}
		return &api.ExecuteResponse{Result: string(resultBytes)}, nil
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 1. 读取请求
	decoder := json.NewDecoder(conn)
	var req api.ExecuteRequest
	if err := decoder.Decode(&req); err != nil {
		if err != io.EOF {
			log.Printf("Decode error: %v", err)
		}
		return
	}

	// 2. 执行请求
	resp, err := Execute(&req)
	if err != nil {
		resp = &api.ExecuteResponse{Error: err.Error()}
	}

	// 3. 发送响应
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(resp); err != nil {
		log.Printf("Encode error: %v", err)
	}
}

func main() {
	// 在 Occlum 环境下，通常由 Init-RA 机制自动处理证明。
	// 这里我们使用标准 TLS。如果需要双向验证或特定证书，可以从环境变量加载路径。

	certFile := os.Getenv("TLS_CERT_PATH")
	keyFile := os.Getenv("TLS_KEY_PATH")

	var tlsConfig *tls.Config
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("加载证书失败: %v", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	} else {
		log.Printf("未提供证书路径 (TLS_CERT_PATH, TLS_KEY_PATH)，将以降级模式或 Init-RA 预期模式运行")
	}

	// 启动监听
	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer ln.Close()

	var tlsLn net.Listener
	if tlsConfig != nil && len(tlsConfig.Certificates) > 0 {
		tlsLn = tls.NewListener(ln, tlsConfig)
		log.Printf("TEE Service running on :%s (TLS enabled)\n", port)
	} else {
		log.Printf("警告：使用未加密连接运行（或等待外界 TLS 终止）")
		tlsLn = ln
		log.Printf("TEE Service running on :%s (Insecure)\n", port)
	}
	defer tlsLn.Close()

	for {
		conn, err := tlsLn.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
