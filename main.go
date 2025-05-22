package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting Chart System...")
	
	// TODO: 初始化配置
	// TODO: 设置数据库连接
	// TODO: 注册路由
	
	// 启动HTTP服务
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
