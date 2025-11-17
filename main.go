package main

import (
	"blog-gin-project/api"
	"blog-gin-project/config"
	"blog-gin-project/database"
	"log"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	database.InitDB()

	// 设置路由
	r := api.SetupRouter()

	// 启动服务器
	log.Printf("服务器将在 %s 端口启动", cfg.ServerPort)
	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
