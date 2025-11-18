package main

import (
	"blog-gin-project/database"
	"blog-gin-project/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()
	defer func() {
		sqlDB, _ := database.DB.DB()
		sqlDB.Close()
	}()

	// 设置Gin模式
	gin.SetMode(gin.DebugMode)

	// 创建Gin引擎
	r := gin.Default()

	// 加载模板和设置路由
	routes.LoadTemplates(r)
	routes.SetupWebRoutes(r)

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "页面未找到",
			"Year":  2025,
		})
	})

	// 启动服务器
	log.Println("服务器将在 :8080 端口启动")
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
