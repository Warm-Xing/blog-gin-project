package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
type Config struct {
	ServerPort     string
	JWTSecret      string
	JWTExpireHours int
}

// LoadConfig 从环境变量或.env文件加载配置
func LoadConfig() Config {
	// 加载.env文件
	godotenv.Load()

	// 获取JWT过期时间，默认为24小时
	jwtExpireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		jwtExpireHours = 24
	}

	return Config{
		ServerPort:     getEnv("SERVER_PORT", ":8080"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-should-be-changed"),
		JWTExpireHours: jwtExpireHours,
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
