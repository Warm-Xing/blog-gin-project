package database

import (
	"blog-gin-project/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	// 使用MySQL数据库
	DB, err = gorm.Open(mysql.Open("root:Wam1992-7414@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 显示SQL日志
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据表
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		log.Fatalf("数据表迁移失败: %v", err)
	}

	log.Println("数据库连接和迁移成功")
}
